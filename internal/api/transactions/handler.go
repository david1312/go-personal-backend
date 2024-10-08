package transactions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"libra-internal/internal/api/products"
	"libra-internal/internal/api/response"
	"libra-internal/pkg/constants"
	"libra-internal/pkg/crashy"
	"libra-internal/pkg/helper"
	"libra-internal/repository/repo_master_data"
	"libra-internal/repository/repo_products"
	"libra-internal/repository/repo_transactions"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	localMdl "libra-internal/internal/api/middleware"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

type TransactionsHandler struct {
	db           *sqlx.DB
	prodRepo     repo_products.ProductsRepository
	mdRepo       repo_master_data.MasterDataRepository
	trRepo       repo_transactions.TransactionsRepositoy
	baseAssetUrl string
	client       *http.Client
	MidtransConfig
	FCMConfig
}

func NewTransactionsHandler(db *sqlx.DB, pr repo_products.ProductsRepository, md repo_master_data.MasterDataRepository, tr repo_transactions.TransactionsRepositoy, baseAssetUrl string, cl *http.Client, midtransCfg MidtransConfig, fcmConfig FCMConfig) *TransactionsHandler {
	return &TransactionsHandler{db: db, prodRepo: pr, baseAssetUrl: baseAssetUrl, mdRepo: md, trRepo: tr, MidtransConfig: midtransCfg, client: cl, FCMConfig: fcmConfig}
}

func (tr *TransactionsHandler) SubmitTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		p                       SubmitTransactionsRequest
		ctx                     = r.Context()
		authData                = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		today                   = time.Now()
		totalBayar              = 0
		transferResponse        = TransferBNIResponse{}
		transferResponsePermata = TransferPermataResponse{}
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	checkTime, _ := time.Parse("2006-01-02", p.ScheduleDate)
	if today.Truncate(24 * time.Hour).After(checkTime.Truncate(24 * time.Hour)) {

		response.Nay(w, r, crashy.New(errors.New(crashy.ErrBackwardDate), crashy.ErrCode(crashy.ErrBackwardDate), crashy.Message(crashy.ErrBackwardDate)), http.StatusBadRequest)
		return
	}

	if !helper.ValidateScheduleTime(p.ScheduleTime) {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidScheduleTime), crashy.ErrCode(crashy.ErrInvalidScheduleTime), crashy.Message(crashy.ErrInvalidScheduleTime)), http.StatusBadRequest)
		return
	}

	lastTransId, errCode, err := tr.trRepo.GetLastTransactionId(ctx)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	custId, errCode, err := tr.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	tranType := p.TranType
	if len(p.TranType) == 0 {
		tranType = constants.TrTypeBooking
	}

	newTransId := helper.GenerateTransactionId(lastTransId, today.Format("20060102"))

	tempListProduct := []repo_transactions.Product{}
	for _, v := range p.ListProduct {
		totalBayar += int(v.Price)
		tempListProduct = append(tempListProduct, repo_transactions.Product{
			ProductId: v.ProductId,
			Qty:       v.Qty,
			Price:     v.Price,
			Total:     v.Price * float64(v.Qty),
		})
	}

	errCode, err = tr.trRepo.SubmitTransaction(ctx, repo_transactions.SubmitTransactionsParam{
		NoFaktur:      newTransId,
		ScheduleDate:  p.ScheduleDate,
		ScheduleTime:  p.ScheduleTime,
		IdOutlet:      p.IdOutlet,
		TranType:      tranType,
		PaymentMethod: p.PaymentMethod,
		Notes:         p.Notes,
		Source:        constants.TransactionSourceApp,
		CustomerId:    custId,
		ListProduct:   tempListProduct,
	})
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	//testing payment
	// var payload struct
	if p.PaymentMethod == constants.COD {
		response.Yay(w, r, SubmitTransactionResponse{
			Status:    constants.StatusSuccess,
			InvoiceId: newTransId,
		}, http.StatusOK)
		return
	}
	payload := &TransferBNIRequest{
		PaymentType: constants.PaymentBankTransfer,
		TransactionDetailsData: TransactionDetails{
			OrderId:     newTransId,
			GrossAmount: fmt.Sprintf("%v", totalBayar),
		},
		BankTransferData: BankTransfer{
			Bank: helper.MappingBankNameRequestMidtrans(p.PaymentMethod),
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), crashy.Message(crashy.ErrRequestMidtrans)), http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", tr.MidtransConfig.BaseUrl, tr.MidtransConfig.ChargeUrl), bytes.NewBuffer(b)) //todo get from config
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), crashy.Message(crashy.ErrRequestMidtrans)), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+tr.MidtransConfig.AuthKey)

	res, err := tr.client.Do(req)

	if err != nil {
		if os.IsTimeout(err) {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), crashy.Message(crashy.ErrRequestMidtrans)), http.StatusInternalServerError)
			return
		}
		return
	}
	if res != nil && res.Body != nil {
		defer func(c io.Closer) {
			err = c.Close()
		}(res.Body)
	}
	j := map[string]interface{}{}
	err = json.NewDecoder(res.Body).Decode(&j)
	if err != nil {
		panic(err)
	}
	if j["status_code"] != "201" {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), fmt.Sprintf("%v: %v", crashy.Message(crashy.ErrRequestMidtrans), j["status_message"])), http.StatusInternalServerError)
		return
	}
	if j["status_code"] != "200" {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), fmt.Sprintf("%v: %v", crashy.Message(crashy.ErrRequestMidtrans), j["status_message"])), http.StatusInternalServerError)
		return
	}

	if p.PaymentMethod == constants.TF_BNI || p.PaymentMethod == constants.TF_BRI {
		_ = json.NewDecoder(res.Body).Decode(&transferResponse)

		_, err = tr.trRepo.UpdateInvoiceVA(ctx, newTransId, transferResponse.VirtualAccountData[0].VaNumber)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), crashy.Message(crashy.ErrRequestMidtrans)), http.StatusInternalServerError)
			return
		}
	}
	if p.PaymentMethod == constants.TF_PERMATA {
		_ = json.NewDecoder(res.Body).Decode(&transferResponsePermata)

		_, err = tr.trRepo.UpdateInvoiceVA(ctx, newTransId, transferResponsePermata.PermataVANumber)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), crashy.Message(crashy.ErrRequestMidtrans)), http.StatusInternalServerError)
			return
		}
	}

	response.Yay(w, r, SubmitTransactionResponse{
		Status:    constants.StatusSuccess,
		InvoiceId: newTransId,
	}, http.StatusOK)

}

func (tr *TransactionsHandler) InquirySchedule(w http.ResponseWriter, r *http.Request) {
	var (
		ctx          = r.Context()
		today        = time.Now()
		end          = today.AddDate(0, 0, 13)
		listSchedule = []InquiryScheduleResponse{}
	)
	data, errCode, err := tr.trRepo.InquirySchedule(ctx, today.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	mappedData := make(map[string]ScheduleCount)
	for _, m := range data {

		mappedData[m.ScheduleDate.Format("2006-01-02")+"-"+m.ScheduleTime] = ScheduleCount{
			ScheduleTime: m.ScheduleTime,
			OrderCount:   m.OrderCount,
		}
	}

	mappedDataResponse := make(map[string][]ScheduleCount)
	for d := today; !d.After(end); d = d.AddDate(0, 0, 1) {
		for _, v := range constants.ScheduleTime {
			orderCount := 0

			if val, ok := mappedData[d.Format("2006-01-02")+"-"+v]; ok {
				orderCount = val.OrderCount
			}

			available := true
			if orderCount >= constants.LimitOrderInHour {
				available = false
			}

			mappedDataResponse[d.Format("2006-01-02")] = append(mappedDataResponse[d.Format("2006-01-02")], ScheduleCount{
				ScheduleTime: v,
				OrderCount:   orderCount,
				IsAvailable:  available,
			})
		}

	}
	keys := make([]string, 0, len(mappedDataResponse))
	for k := range mappedDataResponse {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tempListSchedule := []ScheduleCount{}
		for _, val := range mappedDataResponse[k] {
			tempListSchedule = append(tempListSchedule, ScheduleCount{
				ScheduleTime: val.ScheduleTime,
				OrderCount:   val.OrderCount,
				IsAvailable:  val.IsAvailable,
			})
		}

		listSchedule = append(listSchedule, InquiryScheduleResponse{
			ScheduleDate: k,
			ScheduleList: tempListSchedule,
		},
		)
	}
	response.Yay(w, r, listSchedule, http.StatusOK)
}

func (tr *TransactionsHandler) GetHistoryTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		ctx                   = r.Context()
		fp                    GetListTransactionRequest
		authData              = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		listTransactionRes    = []TransactionsResponse{}
		listProductByInvoices = []repo_transactions.ProductsData{}
	)

	if err := render.Bind(r, &fp); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	limit := fp.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}
	page := fp.Page
	if page < 1 {
		page = 1
	}

	custId, errCode, err := tr.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	listTransaction, totalData, listInvoiceId, errCode, err := tr.trRepo.GetHistoryTransaction(ctx, repo_transactions.GetListTransactionsParam{
		Limit:              limit,
		Page:               page,
		StatusTransactions: fp.TransStatus,
		CustomerId:         custId,
	})
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//get list product for each invoice
	if len(listInvoiceId) > 0 {
		listProductByInvoices, errCode, err = tr.trRepo.GetProductByInvoices(ctx, listInvoiceId)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
	}

	mappedProductByInvoice := make(map[string][]ProductsData)
	for _, m := range listProductByInvoices {
		mappedProductByInvoice[m.InvoiceId] = append(mappedProductByInvoice[m.InvoiceId], ProductsData{
			KodePLU:              m.KodePLU,
			NamaBarang:           m.NamaBarang,
			NamaUkuran:           m.NamaUkuran,
			Qty:                  m.Qty,
			HargaSatuan:          m.Harga,
			HargaSatuanFormatted: helper.FormatCurrency(int(m.Harga)),
			HargaTotal:           m.HargaTotal,
			HargaTotalFormatted:  helper.FormatCurrency(int(m.HargaTotal)),
			Deskripsi:            m.Deskripsi,
			DisplayImage:         tr.baseAssetUrl + constants.ProductDir + m.DisplayImage,
			JenisBan:             m.JenisBan,
		})

	}

	for _, v := range listTransaction {
		listTransactionRes = append(listTransactionRes, TransactionsResponse{
			InvoiceId:            v.InvoiceId,
			Status:               v.Status,
			TotalAmount:          v.TotalAmount,
			TotalAmountFormatted: helper.FormatCurrency(int(v.TotalAmount)),
			PaymentMethodDesc:    v.PaymentMethodDesc,
			PaymentMethodIcon:    tr.baseAssetUrl + constants.PaymentMethod + v.PaymentMethodIcon,
			OutletId:             v.OutletId,
			OutletName:           v.OutletName,
			CsNumber:             constants.CSNumber,
			CreatedAt:            v.CreatedAt.Format("02 January 2006"),
			PaymentDue:           v.PaymentDue.Format("02 January 2006 15:04"),
			ListProduct:          mappedProductByInvoice[v.InvoiceId],
			InstallationtTime:    helper.FormatInstallationTime(v.InstallationDate, v.InstallationTime),
		})
	}

	response.Yay(w, r, ListProductsResponse{
		TransactionData: listTransactionRes,
		DataInfo: products.DataInfo{
			CurrentPage: page,
			MaxPage: func() int {
				maxPage := float64(totalData) / float64(limit)
				if helper.IsFloatNoDecimal(maxPage) {
					return int(maxPage)
				}
				return int(maxPage) + 1
			}(),
			Limit:       limit,
			TotalRecord: totalData,
		},
	}, http.StatusOK)

}

func (tr *TransactionsHandler) CallbackPayment(w http.ResponseWriter, r *http.Request) {
	var (
		p             PaymentCallbackRequest
		ctx           = r.Context()
		paymentStatus = constants.DBPaymentSettle
		transStatus   = constants.TranStatusDipasang
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	if p.TransactionStatus != constants.MTransStatusSettlement {
		paymentStatus = constants.DBPaymentNotSettle
		transStatus = constants.TranStatusPembayaran
	}

	errCode, err := tr.trRepo.UpdateInvoiceStatus(ctx, p.OrderId, transStatus, paymentStatus)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//send notif firebase to gcm id
	userData, _, err := tr.trRepo.GetUserFCMToken(ctx, p.OrderId)
	if err != nil {
		log.Printf("there's an error when sending firebase notification err step A: %v", err.Error())
		response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
		return
	}
	if userData.DeviceToken != "" {
		payload := &FCMRequest{
			To: userData.DeviceToken,
			NotificationData: Notification{
				Body:  fmt.Sprintf(constants.PushNotifMsgSuccess, p.OrderId),
				Title: constants.PushNotifTitle,
			},
			FCMData: DataFCM{
				Action:    constants.PushNotifAction,
				InvoiceID: p.OrderId,
			},
		}
		b, err := json.Marshal(payload)
		if err != nil {
			log.Printf("there's an error when sending firebase notification err step B: %v", err.Error())
			response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
			return
		}
		req, err := http.NewRequest(http.MethodPost, tr.FCMConfig.NotifUrl, bytes.NewBuffer(b)) //todo get from config
		if err != nil {
			log.Printf("there's an error when sending firebase notification err step C: %v", err.Error())
			response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
			return
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tr.FCMConfig.ClientKey)

		_, err = tr.client.Do(req)

		if err != nil {
			log.Printf("there's an error when sending firebase notification err step D: %v", err.Error())
			response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
			return
		}

	}
	fmt.Println("temporary log callback accepted ")
	response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
}

func (tr *TransactionsHandler) GetPaymentInstruction(w http.ResponseWriter, r *http.Request) {
	var (
		p   GetPaymentInstructionRequest
		ctx = r.Context()
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	transaction, errCode, err := tr.trRepo.GetInvoiceData(ctx, p.InvoiceId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if transaction.PaymentMethod == "COD" {
		response.Yay(w, r, GetPaymentInstructionResponse{
			PaymentMethodDesc:    transaction.PaymentMethodDesc,
			PaymentMethodIcon:    tr.baseAssetUrl + constants.PaymentMethod + transaction.PaymentMethodIcon,
			TotalAmountFormatted: helper.FormatCurrency(int(transaction.TotalAmount)),
			VirtualAccNumber:     transaction.VirtualAccount,
			AtmTitle:             "Bayar langsung di Outlet / COD",
			AtmInstruction:       constants.CODIntruction,
			IBInstruction:        []string{},
			MBInstuction:         []string{},
		}, http.StatusOK)
		return
	} else if transaction.PaymentMethod == constants.TF_PERMATA {
		response.Yay(w, r, GetPaymentInstructionResponse{
			PaymentMethodDesc:    transaction.PaymentMethodDesc,
			PaymentMethodIcon:    tr.baseAssetUrl + constants.PaymentMethod + transaction.PaymentMethodIcon,
			TotalAmountFormatted: helper.FormatCurrency(int(transaction.TotalAmount)),
			VirtualAccNumber:     transaction.VirtualAccount,
			AtmTitle:             "ATM " + helper.MappingBankName(transaction.PaymentMethod),
			AtmInstruction:       constants.ATMInstructionPermata,
			IBInstruction:        constants.InternetBankingInstructionPermata,
			MBInstuction:         constants.MobileBankingInstructionPermata,
		}, http.StatusOK)
		return
	} else if transaction.PaymentMethod == constants.TF_BRI {
		response.Yay(w, r, GetPaymentInstructionResponse{
			PaymentMethodDesc:    transaction.PaymentMethodDesc,
			PaymentMethodIcon:    tr.baseAssetUrl + constants.PaymentMethod + transaction.PaymentMethodIcon,
			TotalAmountFormatted: helper.FormatCurrency(int(transaction.TotalAmount)),
			VirtualAccNumber:     transaction.VirtualAccount,
			AtmTitle:             "ATM " + helper.MappingBankName(transaction.PaymentMethod),
			AtmInstruction:       constants.ATMInstructionBRI,
			IBInstruction:        constants.InternetBankingInstructionBRI,
			MBInstuction:         constants.MobileBankingInstructionBRI,
		}, http.StatusOK)
		return
	}

	response.Yay(w, r, GetPaymentInstructionResponse{
		PaymentMethodDesc:    transaction.PaymentMethodDesc,
		PaymentMethodIcon:    tr.baseAssetUrl + constants.PaymentMethod + transaction.PaymentMethodIcon,
		TotalAmountFormatted: helper.FormatCurrency(int(transaction.TotalAmount)),
		VirtualAccNumber:     transaction.VirtualAccount,
		AtmTitle:             "ATM " + helper.MappingBankName(transaction.PaymentMethod),
		AtmInstruction:       constants.ATMInstructionBNI,
		IBInstruction:        constants.InternetBankingInstructionBNI,
		MBInstuction:         constants.MobileBankingInstructionBNI,
	}, http.StatusOK)
}

func (tr *TransactionsHandler) GetTransactionDetail(w http.ResponseWriter, r *http.Request) {
	var (
		p                      GetPaymentInstructionRequest
		ctx                    = r.Context()
		listProductByInvoiceId = []ProductsDataPageJadwal{}
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	transaction, errCode, err := tr.trRepo.GetTransactionDetail(ctx, p.InvoiceId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	listProduct, errCode, err := tr.trRepo.GetProductByInvoiceId(ctx, p.InvoiceId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	for _, m := range listProduct {
		listProductByInvoiceId = append(listProductByInvoiceId, ProductsDataPageJadwal{
			KodePLU:              m.KodePLU,
			NamaBarang:           m.NamaBarang,
			NamaUkuran:           m.NamaUkuran,
			Qty:                  m.Qty,
			HargaSatuan:          m.Harga,
			HargaSatuanFormatted: helper.FormatCurrency(int(m.Harga)),
			HargaTotal:           m.HargaTotal,
			HargaTotalFormatted:  helper.FormatCurrency(int(m.HargaTotal)),
			Deskripsi:            m.Deskripsi,
			DisplayImage:         tr.baseAssetUrl + constants.ProductDir + m.DisplayImage,
			JenisBan:             m.JenisBan,
		},
		)
	}

	splitStr := strings.Split(transaction.InstallationTime, ":")
	timeAdded, _ := strconv.Atoi(splitStr[1])
	timeAdded = timeAdded + 15
	rescheduleTime := fmt.Sprintf("%v:%v", splitStr[0], timeAdded)
	//
	date, _ := time.Parse("2006-01-02", transaction.InstallationDate[:10])
	now := time.Now()
	// nowZeroTime
	bannerMsg := constants.BannerMsgToday
	if date.After(now) {
		days := math.Ceil(date.Sub(now).Hours() / 24)
		bannerMsg = fmt.Sprintf(constants.BannerMsgUpcoming, days)
	}
	isEnableReview := false
	if transaction.Status == constants.TransStatusBerhasil {
		isEnableReview = true
	}
	// fmt.Println(transaction.Status)
	// fmt.Println(transaction.InstallationDate +" " + transaction.InstallationTime)

	response.Yay(w, r, GetTransactionsDetailResponse{
		InvoiceId:          transaction.InvoiceId,
		BannerInformation:  bannerMsg,
		InstallationtTime:  helper.FormatInstallationTime(transaction.InstallationDate, transaction.InstallationTime),
		OutletName:         transaction.OutletName,
		CsNumber:           constants.CSNumber,
		OutletAddress:      fmt.Sprintf("%v, %v, %v", transaction.OutletAddress, transaction.OutletDistrict, transaction.OutletCity),
		RescheduleTime:     helper.FormatInstallationTime(transaction.InstallationDate, rescheduleTime),
		IsEnableReview:     isEnableReview,
		IsEnableReschedule: false,
		PaymentMethod:      transaction.PaymentMethod,
		PaymentMethodDesc:  transaction.PaymentMethodDesc,
		PaymentMethodIcon:  tr.baseAssetUrl + constants.PaymentMethod + transaction.PaymentMethodIcon,
		ListProduct:        listProductByInvoiceId,
	}, http.StatusOK)

}

func (tr *TransactionsHandler) GetCountTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	custId, errCode, err := tr.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	res, errCode, err := tr.trRepo.GetCountTransactionData(ctx, custId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, GetSummaryTransactionCountResponse{
		WaitingPayment: res.WaitingPayment,
		WaitingProcess: res.WaitingProcess,
		OnProgress:     res.OnProgress,
		Succedd:        res.Succedd,
	}, http.StatusOK)
}

func (tr *TransactionsHandler) EPUpdateTransactionStatus(w http.ResponseWriter, r *http.Request) {
	var (
		p   TransactionCommonRequest
		ctx = r.Context()
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := tr.trRepo.UpdateTransactionStatus(ctx, p.InvoiceId, p.Status, p.Notes)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if p.Status == constants.TranStatusDipasang {
		//send notif firebase to gcm id
		userData, _, err := tr.trRepo.GetUserFCMToken(ctx, p.InvoiceId)
		if err != nil {
			log.Printf("there's an error when sending firebase notification err step A: %v", err.Error())
			response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
			return
		}
		if userData.DeviceToken != "" {
			payload := &FCMRequest{
				To: userData.DeviceToken,
				NotificationData: Notification{
					Body:  fmt.Sprintf(constants.PushNotifMsgSuccess, p.InvoiceId),
					Title: constants.PushNotifTitle,
				},
				FCMData: DataFCM{
					Action:    constants.PushNotifAction,
					InvoiceID: p.InvoiceId,
				},
			}
			b, err := json.Marshal(payload)
			if err != nil {
				log.Printf("there's an error when sending firebase notification err step B: %v", err.Error())
				response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
				return
			}
			req, err := http.NewRequest(http.MethodPost, tr.FCMConfig.NotifUrl, bytes.NewBuffer(b)) //todo get from config
			if err != nil {
				log.Printf("there's an error when sending firebase notification err step C: %v", err.Error())
				response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
				return
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tr.FCMConfig.NotifUrl)

			_, err = tr.client.Do(req)

			if err != nil {
				log.Printf("there's an error when sending firebase notification err step D: %v", err.Error())
				response.Yay(w, r, constants.StatusSuccess, http.StatusOK)
				return
			}

		}

	}

	response.Yay(w, r, constants.StatusSuccess, http.StatusOK)

}

func (tr *TransactionsHandler) EPMerchantGetTransactionDetail(w http.ResponseWriter, r *http.Request) {
	var (
		p                      GetPaymentInstructionRequest
		ctx                    = r.Context()
		listProductByInvoiceId = []ProductsDataPageJadwal{}
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	transaction, errCode, err := tr.trRepo.GetTransactionDetail(ctx, p.InvoiceId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	listProduct, errCode, err := tr.trRepo.GetProductByInvoiceId(ctx, p.InvoiceId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	for _, m := range listProduct {
		listProductByInvoiceId = append(listProductByInvoiceId, ProductsDataPageJadwal{
			KodePLU:              m.KodePLU,
			NamaBarang:           m.NamaBarang,
			NamaUkuran:           m.NamaUkuran,
			Qty:                  m.Qty,
			HargaSatuan:          m.Harga,
			HargaSatuanFormatted: helper.FormatCurrency(int(m.Harga)),
			HargaTotal:           m.HargaTotal,
			HargaTotalFormatted:  helper.FormatCurrency(int(m.HargaTotal)),
			Deskripsi:            m.Deskripsi,
			DisplayImage:         tr.baseAssetUrl + constants.ProductDir + m.DisplayImage,
			JenisBan:             m.JenisBan,
		},
		)
	}

	splitStr := strings.Split(transaction.InstallationTime, ":")
	timeAdded, _ := strconv.Atoi(splitStr[1])
	timeAdded = timeAdded + 15
	rescheduleTime := fmt.Sprintf("%v:%v", splitStr[0], timeAdded)
	//
	date, _ := time.Parse("2006-01-02", transaction.InstallationDate[:10])
	now := time.Now()
	// nowZeroTime
	bannerMsg := constants.BannerMsgToday
	if date.After(now) {
		days := math.Ceil(date.Sub(now).Hours() / 24)
		bannerMsg = fmt.Sprintf(constants.BannerMsgUpcoming, days)
	}
	// isEnableReview := false
	// if transaction.Status == constants.TransStatusBerhasil {
	// 	isEnableReview = true
	// }
	// fmt.Println(transaction.Status)
	// fmt.Println(transaction.InstallationDate +" " + transaction.InstallationTime)

	response.Yay(w, r, GetTransactionsDetaiMerchantlResponse{
		InvoiceId:          transaction.InvoiceId,
		Status:             transaction.Status,
		BannerInformation:  bannerMsg,
		InstallationtTime:  helper.FormatInstallationTime(transaction.InstallationDate, transaction.InstallationTime),
		OutletName:         transaction.OutletName,
		CsNumber:           constants.CSNumber,
		OutletAddress:      fmt.Sprintf("%v, %v, %v", transaction.OutletAddress, transaction.OutletDistrict, transaction.OutletCity),
		RescheduleTime:     helper.FormatInstallationTime(transaction.InstallationDate, rescheduleTime),
		IsEnableReview:     true,
		IsEnableReschedule: false,
		PaymentMethod:      transaction.PaymentMethod,
		PaymentMethodDesc:  transaction.PaymentMethodDesc,
		PaymentMethodIcon:  tr.baseAssetUrl + constants.PaymentMethod + transaction.PaymentMethodIcon,
		ListProduct:        listProductByInvoiceId,
		CustomerName:       transaction.CustomerName,
		CustomerPhone:      transaction.CustomerPhone,
		CustomerEmail:      transaction.CustomerEmail,
		Notes:              transaction.Notes,
	}, http.StatusOK)

}

func (tr *TransactionsHandler) EPMerchantGetHistoryTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		ctx                   = r.Context()
		fp                    GetListTransactionRequest
		listTransactionRes    = []TransactionsResponseMerchant{}
		listProductByInvoices = []repo_transactions.ProductsData{}
	)

	if err := render.Bind(r, &fp); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	limit := fp.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}
	page := fp.Page
	if page < 1 {
		page = 1
	}

	listTransaction, totalData, listInvoiceId, errCode, err := tr.trRepo.GetHistoryTransactionMerchant(ctx, repo_transactions.GetListTransactionsParam{
		Limit:              limit,
		Page:               page,
		StatusTransactions: fp.TransStatus,
	})
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//get list product for each invoice
	if len(listInvoiceId) > 0 {
		listProductByInvoices, errCode, err = tr.trRepo.GetProductByInvoices(ctx, listInvoiceId)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
			return
		}
	}

	mappedProductByInvoice := make(map[string][]ProductsData)
	for _, m := range listProductByInvoices {
		mappedProductByInvoice[m.InvoiceId] = append(mappedProductByInvoice[m.InvoiceId], ProductsData{
			KodePLU:              m.KodePLU,
			NamaBarang:           m.NamaBarang,
			NamaUkuran:           m.NamaUkuran,
			Qty:                  m.Qty,
			HargaSatuan:          m.Harga,
			HargaSatuanFormatted: helper.FormatCurrency(int(m.Harga)),
			HargaTotal:           m.HargaTotal,
			HargaTotalFormatted:  helper.FormatCurrency(int(m.HargaTotal)),
			Deskripsi:            m.Deskripsi,
			DisplayImage:         tr.baseAssetUrl + constants.ProductDir + m.DisplayImage,
			JenisBan:             m.JenisBan,
		})

	}

	for _, v := range listTransaction {
		listTransactionRes = append(listTransactionRes, TransactionsResponseMerchant{
			InvoiceId:            v.InvoiceId,
			Status:               v.Status,
			TotalAmount:          v.TotalAmount,
			TotalAmountFormatted: helper.FormatCurrency(int(v.TotalAmount)),
			PaymentMethodDesc:    v.PaymentMethodDesc,
			PaymentMethodIcon:    tr.baseAssetUrl + constants.PaymentMethod + v.PaymentMethodIcon,
			OutletId:             v.OutletId,
			OutletName:           v.OutletName,
			CsNumber:             constants.CSNumber,
			CreatedAt:            v.CreatedAt.Format("02 January 2006"),
			PaymentDue:           v.PaymentDue.Format("02 January 2006 15:04"),
			ListProduct:          mappedProductByInvoice[v.InvoiceId],
			CustomerName:         v.CustomerName,
			CustomerPhone:        v.CustomerPhone,
			CustomerEmail:        v.CustomerEmail,
		})
	}

	response.Yay(w, r, ListProductsResponseMerchant{
		TransactionData: listTransactionRes,
		DataInfo: products.DataInfo{
			CurrentPage: page,
			MaxPage: func() int {
				maxPage := float64(totalData) / float64(limit)
				if helper.IsFloatNoDecimal(maxPage) {
					return int(maxPage)
				}
				return int(maxPage) + 1
			}(),
			Limit:       limit,
			TotalRecord: totalData,
		},
	}, http.StatusOK)

}
