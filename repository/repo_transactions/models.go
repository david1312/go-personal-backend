package repo_transactions

import "time"

type SubmitTransactionsParam struct {
	NoFaktur      string
	ScheduleDate  string
	ScheduleTime  string
	IdOutlet      int
	TranType      string
	PaymentMethod string
	Notes         string
	Source        string
	CustomerId    int
	ListProduct   []Product
}

type Product struct {
	ProductId int
	Qty       int
	Price     float64
	Total     float64
}

type ScheduleCount struct {
	ScheduleDate time.Time
	ScheduleTime string
	OrderCount   int
}
