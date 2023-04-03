package reports

import (
	"fmt"
	"libra-internal/internal/api/middleware"
	"libra-internal/internal/api/response"
	"libra-internal/pkg/crashy"
	"libra-internal/pkg/helper"
	"libra-internal/repository/repo_reports"
	"net/http"
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
		ctx      = r.Context()
		authData = ctx.Value(middleware.CtxKey).(middleware.MerchantToken)
	)

	fileName, errCode, err := helper.UploadSingleFile(r, "report", "/files", "/reports", "sales-*.xlsx", 10)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusBadRequest)
		return
	}

	fmt.Println(fileName)

	fmt.Println(authData.Username)

}
