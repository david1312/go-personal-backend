package repo_transactions

import "context"

type TransactionsRepositoy interface {
	GetLastTransactionId(ctx context.Context) (res, errCode string, err error)
	SubmitTransaction(ctx context.Context, fp SubmitTransactionsParam) (errCode string, err error)
	InquirySchedule(ctx context.Context, startDate, endDate string) (res []ScheduleCount, errCode string, err error)
}
