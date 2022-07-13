package transactions

import (
	"errors"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/repository/repo_master_data"
	"semesta-ban/repository/repo_products"
	"semesta-ban/repository/repo_transactions"
	"sort"
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
}

func NewTransactionsHandler(db *sqlx.DB, pr repo_products.ProductsRepository, md repo_master_data.MasterDataRepository, tr repo_transactions.TransactionsRepositoy, baseAssetUrl string) *TransactionsHandler {
	return &TransactionsHandler{db: db, prodRepo: pr, baseAssetUrl: baseAssetUrl, mdRepo: md, trRepo: tr}
}

func (tr *TransactionsHandler) SubmitTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		p        SubmitTransactionsRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		today    = time.Now()
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

	response.Yay(w, r, "success", http.StatusOK)

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
