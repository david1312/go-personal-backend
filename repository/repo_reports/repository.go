package repo_reports

import (
	"context"
	"libra-internal/internal/models"
	"net/http"
)

type ReportsRepository interface {
	SyncUpSales(ctx context.Context, fileName, dir string) (err error)
	UpdateNetProfit(ctx context.Context, limit bool) (err error)
	GetAllSalesReport(ctx context.Context, params models.GetAllSalesRequest) (res []SalesModel, pageData models.Pagination, summary models.SummarySales, errCode string, err error)
	GetSalesByInvoice(ctx context.Context, noPesanan string, client *http.Client) (res models.ApiResponseSalesDetail, errCode string, err error)
	GetAllSalesMinusReport(ctx context.Context, params models.GetAllSalesRequest) (res []SalesModel, pageData models.Pagination, summary models.SummarySales, errCode string, err error)
}
