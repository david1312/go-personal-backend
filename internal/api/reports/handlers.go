package reports

import (
	"fmt"
	"libra-internal/internal/api/response"
	"libra-internal/internal/models"
	"libra-internal/pkg/constants"
	"libra-internal/pkg/crashy"
	"libra-internal/pkg/helper"
	"libra-internal/repository/repo_reports"
	"net/http"

	"github.com/go-chi/render"
)

type ReportsHandler struct {
	reportsRepo repo_reports.ReportsRepository
}

//todo REMEMBER 30 May gmail tidak support lagi less secure app find solution

func NewReportsHandler(rr repo_reports.ReportsRepository) *ReportsHandler {
	return &ReportsHandler{reportsRepo: rr}
}

func (rep *ReportsHandler) SyncSales(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		// authData = ctx.Value(middleware.CtxKey).(middleware.MerchantToken)
	)

	fileName, errCode, err := helper.UploadSingleFile(r, "report", constants.DIR_FILES, constants.DIR_REPORT_SALES, constants.FORMAT_EXCEL, constants.MAX_COMMON_SIZE)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	//todo insert into history if success sync up data
	err = rep.reportsRepo.SyncUpSales(ctx, fileName, fmt.Sprintf("%v%v", constants.DIR_FILES, constants.DIR_REPORT_SALES))
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeDataWrite, crashy.Message(crashy.ErrCodeDataWrite)), http.StatusInternalServerError)
		return
	}

	// fmt.Println(fileName)

	// fmt.Println(authData.Username)
	response.Yay(w, r, "success", http.StatusOK)

}

func (rep *ReportsHandler) SalesCalculateProfit(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		// authData = ctx.Value(middleware.CtxKey).(middleware.MerchantToken)
	)

	err := rep.reportsRepo.UpdateNetProfit(ctx, true)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeDataWrite, crashy.Message(crashy.ErrCodeDataWrite)), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (rep *ReportsHandler) EPGetSalesReport(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		fp  models.GetAllSalesRequest
		// authData  = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		salesResponseList = []models.SalesResponse{}
	)

	if err := render.Bind(r, &fp); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	limit := fp.Limit
	if limit < 1 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	page := fp.Page
	if page < 1 {
		page = 1
	}

	salesData, paginationData, summarySales, errCode, err := rep.reportsRepo.GetAllSalesReport(ctx, models.GetAllSalesRequest{
		Limit:     limit,
		Page:      page,
		StartDate: fp.StartDate,
		EndDate:   fp.EndDate,
		NoPesanan: fp.NoPesanan,
	})
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	for _, val := range salesData {
		salesResponseList = append(salesResponseList, models.SalesResponse{
			ID:                  val.ID,
			Tanggal:             val.Tanggal,
			NoPesanan:           val.NoPesanan,
			Status:              val.Status,
			Channel:             val.Channel,
			NettSales:           val.NettSales,
			GrossProfit:         val.GrossProfit,
			PotonganMarketplace: val.PotonganMarketPlace,
			NetProfit:           val.NetProfit,
		})
	}

	response.Yay(w, r, models.ApiResponseSales{
		PaginationData: paginationData,
		SummaryData:    summarySales,
		SalesList:      salesResponseList,
	}, http.StatusOK)
}
