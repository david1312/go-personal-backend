package transactions

import (
	"net/http"
	"semesta-ban/internal/api/products"

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

type GetListTransactionRequest struct {
	Limit       int      `json:"limit"`
	Page        int      `json:"page"`
	TransStatus []string `json:"trans_status"`
}

func (m *GetListTransactionRequest) Bind(r *http.Request) error {
	return m.ValidateGetListTransactionRequest()
}

func (m *GetListTransactionRequest) ValidateGetListTransactionRequest() error {
	return validation.ValidateStruct(m)
}

type TransactionsResponse struct {
	InvoiceId            string         `json:"invoice_id"`
	Status               string         `json:"status"`
	TotalAmount          float64        `json:"total_amount"`
	TotalAmountFormatted string         `json:"total_amount_formatted"`
	PaymentMethodDesc    string         `json:"payment_method_desc"`
	PaymentMethodIcon    string         `json:"payment_method_icon"`
	PaymentDue           string         `json:"payment_due"`
	CreatedAt            string         `json:"created_at"`
	ListProduct          []ProductsData `json:"list_product"`
}

type ProductsData struct {
	KodePLU              int32   `json:"id"`
	NamaBarang           string  `json:"nama_barang"`
	NamaUkuran           string  `json:"ukuran"`
	Qty                  int     `json:"qty"`
	HargaSatuan          float64 `json:"harga_satuan"`
	HargaSatuanFormatted string  `json:"harga_satuan_formatted"`
	HargaTotal           float64 `json:"harga_total"`
	HargaTotalFormatted  string  `json:"harga_total_formatted"`
	Deskripsi            string  `json:"deskripsi"`
	DisplayImage         string  `json:"display_image"`
}

type ListProductsResponse struct {
	DataInfo        products.DataInfo      `json:"info"`
	TransactionData []TransactionsResponse `json:"data"`
}

type MidtransConfig struct {
	MerchantId string
	ClientKey  string
	ServerKey  string
}
