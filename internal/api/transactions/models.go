package transactions

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

type SubmitTransactionsRequest struct {
	ScheduleDate  string    `json:"schedule_date"`
	ScheduleTime  string    `json:"schedule_time"`
	IdOutlet      int       `json:"id_outlet"`
	TranType      string    `json:"tran_type"`
	PaymentMethod string    `json:"payment_method"`
	Notes         string    `json:"notes"`
	ListProduct   []Product `json:"list_product"`
}

type Product struct {
	ProductId int     `json:"product_id"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
}

func (m *SubmitTransactionsRequest) Bind(r *http.Request) error {
	return m.ValidateSubmitTransactionsRequest()
}

func (m *SubmitTransactionsRequest) ValidateSubmitTransactionsRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.ScheduleDate, validation.Required),
		validation.Field(&m.ScheduleTime, validation.Required),
		validation.Field(&m.IdOutlet, validation.Required),
		validation.Field(&m.PaymentMethod, validation.Required),
		validation.Field(&m.ListProduct, validation.Required),
	)
}

type ScheduleCount struct {
	ScheduleTime string `json:"schedule_time"`
	OrderCount   int    `json:"order_count"`
	IsAvailable  bool   `json:"is_available"`
}

type InquiryScheduleResponse struct {
	ScheduleDate string          `json:"schedule_date"`
	ScheduleList []ScheduleCount `json:"list_schedule"`
}

type ScheduleCountOnly struct {
	OrderCount  int  `json:"order_count"`
	IsAvailable bool `json:"is_available"`
}

type ScheduleCountMap struct {
	ScheduleTime          string
	ScheduleCountOnlyList []ScheduleCountOnly
}
