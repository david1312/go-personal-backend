package repo_reports

import (
	"context"
	"libra-internal/internal/models"
)

type ReportsRepository interface {
	SyncUpSales(ctx context.Context, fileName, dir string) (err error)
	UpdateNetProfit(ctx context.Context, limit bool) (err error)
	GetDetailInvoice(ctx context.Context, noPesanan string) (err error)
	GetAllSalesReport(ctx context.Context, params models.GetAllSalesRequest) (res []SalesModel, pageData models.Pagination, summary models.SummarySales, errCode string, err error)
	GetSalesByInvoice(ctx context.Context, noPesanan string) (res models.ApiResponseSalesDetail, errCode string, err error)
}
