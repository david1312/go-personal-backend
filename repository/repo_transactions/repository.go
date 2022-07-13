package repo_transactions

import "context"

type TransactionsRepositoy interface {
	GetLastTransactionId(ctx context.Context) (res, errCode string, err error)
	SubmitTransaction(ctx context.Context, fp SubmitTransactionsParam) (errCode string, err error)
	InquirySchedule(ctx context.Context, startDate, endDate string) (res []ScheduleCount, errCode string, err error)
	GetHistoryTransaction(ctx context.Context, fp GetListTransactionsParam) (res []Transactions, totalData int, listInvoice []string, errCode string, err error)
	GetProductByInvoices(ctx context.Context, listInvoiceId []string) (res []ProductsData, errCode string, err error)
}
