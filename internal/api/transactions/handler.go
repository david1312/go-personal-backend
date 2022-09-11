package transactions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"semesta-ban/internal/api/products"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/repository/repo_master_data"
	"semesta-ban/repository/repo_products"
	"semesta-ban/repository/repo_transactions"
	"sort"
	"strconv"
	"strings"
	"time"

	localMdl "semesta-ban/internal/api/middleware"

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
}

func NewTransactionsHandler(db *sqlx.DB, pr repo_products.ProductsRepository, md repo_master_data.MasterDataRepository, tr repo_transactions.TransactionsRepositoy, baseAssetUrl string, cl *http.Client, midtransCfg MidtransConfig) *TransactionsHandler {
	return &TransactionsHandler{db: db, prodRepo: pr, baseAssetUrl: baseAssetUrl, mdRepo: md, trRepo: tr, MidtransConfig: midtransCfg, client: cl}
}

func (tr *TransactionsHandler) SubmitTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		p                SubmitTransactionsRequest
		ctx              = r.Context()
		authData         = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		today            = time.Now()
		totalBayar       = 0
		transferResponse = TransferBNIResponse{}
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
		Source:        "APP",
		CustomerId:    custId,
		ListProduct:   tempListProduct,
	})
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	//testing payment
	payload := &TransferBNIRequest{
		PaymentType: constants.PaymentBankTransfer,
		TransactionDetailsData: TransactionDetails{
			OrderId:     newTransId,
			GrossAmount: fmt.Sprintf("%v", totalBayar),
		},
		BankTransferData: BankTransfer{
			Bank: constants.BankBNI,
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), crashy.Message(crashy.ErrRequestMidtrans)), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.sandbox.midtrans.com/v2/charge", bytes.NewBuffer(b)) //todo get from config
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

	_ = json.NewDecoder(res.Body).Decode(&transferResponse)
	// fmt.Println(transferResponse.VirtualAccountData[0].VaNumber)

	errCode, err = tr.trRepo.UpdateInvoiceVA(ctx, newTransId, transferResponse.VirtualAccountData[0].VaNumber)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(crashy.ErrRequestMidtrans), crashy.Message(crashy.ErrRequestMidtrans)), http.StatusInternalServerError)
		return
	}
	//end test
	//todo update data to net table payment from this response

	response.Yay(w, r, SubmitTransactionResponse{
		Status:    "success",
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
		transStatus   = "Menunggu Dipasang"
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	if p.TransactionStatus != constants.MTransStatusSettlement {
		paymentStatus = constants.DBPaymentNotSettle
		transStatus = "Menunggu Pembayaran"
	}
	errCode, err := tr.trRepo.UpdateInvoiceStatus(ctx, p.OrderId, transStatus, paymentStatus)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//send notif firebase to gcm id
	fmt.Println("temporary log callback accepted ")
	response.Yay(w, r, "success", http.StatusOK)
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
	bannerMsg := ""
	if helper.DateEqual(date, now) {
		bannerMsg = "Ban pilihan mu akan dipasang hari ini"
	} else if date.After(now) {
		days := math.Ceil(date.Sub(now).Hours() / 24)
		bannerMsg = fmt.Sprintf("Ban pilihan mu akan dipasang dalam %v hari", days)
	}
	// isEnableReview := false
	// if transaction.Status == constants.TransStatusBerhasil {
	// 	isEnableReview = true
	// }
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
		IsEnableReview:     true,
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
