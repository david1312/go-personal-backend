package repo_reports

import "context"

type ReportsRepository interface {
	SyncUpSales(ctx context.Context, fileName string) (err error)
}
