package repo_transactions

import "context"

type TransactionsRepositoy interface {
	GetLastTransactionId(ctx context.Context) (res, errCode string, err error)
	SubmitTransaction(ctx context.Context, fp SubmitTransactionsParam) (errCode string, err error)
	InquirySchedule(ctx context.Context, startDate, endDate string) (res []ScheduleCount, errCode string, err error)
	GetHistoryTransaction(ctx context.Context, fp GetListTransactionsParam) (res []Transactions, totalData int, listInvoice []string, errCode string, err error)
	GetProductByInvoices(ctx context.Context, listInvoiceId []string) (res []ProductsData, errCode string, err error)
	UpdateInvoiceVA(ctx context.Context, invoiceId, virtualAcc string) (errCode string, err error)
	UpdateInvoiceStatus(ctx context.Context, invoiceId, transStatus, paymentStatus string) (errCode string, err error)
	GetInvoiceData(ctx context.Context, invoiceId string) (res Transactions, errCode string, err error)
	GetProductByInvoiceId(ctx context.Context, invoiceId string) (res []ProductsData, errCode string, err error)
	GetTransactionDetail(ctx context.Context, invoiceId string) (res GetTransactionsDetailData, errCode string, err error)
	//
	GetCountTransactionData(ctx context.Context, custId int) (res GetSummaryTransactionCount, errCode string, err error)
	UpdateTransactionStatus(ctx context.Context, invoiceId, status, notes string) (errCode string, err error)
	GetHistoryTransactionMerchant(ctx context.Context, fp GetListTransactionsParam) (res []Transactions, totalData int, listInvoice []string, errCode string, err error)
}
