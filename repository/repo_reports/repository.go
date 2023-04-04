package repo_reports

import "context"

type ReportsRepository interface {
	SyncUpSales(ctx context.Context, fileName, dir string) (err error)
	UpdateNetProfit(ctx context.Context) (err error)
	GetDetailInvoice(ctx context.Context, noPesanan string) (err error)
}
