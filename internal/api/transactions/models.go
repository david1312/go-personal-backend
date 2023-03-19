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
	OutletId             int            `json:"outlet_id"`
	OutletName           string         `json:"outlet_name"`
	CsNumber             string         `json:"outlet_cs_number"`
	CreatedAt            string         `json:"created_at"`
	ListProduct          []ProductsData `json:"list_product"`
	InstallationtTime    string         `json:"installation_time"`
}

type TransactionsResponseMerchant struct {
	InvoiceId            string         `json:"invoice_id"`
	Status               string         `json:"status"`
	TotalAmount          float64        `json:"total_amount"`
	TotalAmountFormatted string         `json:"total_amount_formatted"`
	PaymentMethodDesc    string         `json:"payment_method_desc"`
	PaymentMethodIcon    string         `json:"payment_method_icon"`
	PaymentDue           string         `json:"payment_due"`
	OutletId             int            `json:"outlet_id"`
	OutletName           string         `json:"outlet_name"`
	CsNumber             string         `json:"outlet_cs_number"`
	CreatedAt            string         `json:"created_at"`
	ListProduct          []ProductsData `json:"list_product"`
	CustomerName         string         `json:"customer_name"`
	CustomerPhone        string         `json:"customer_phone"`
	CustomerEmail        string         `json:"customer_email"`
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
	JenisBan             string  `json:"jenis_ban"`
}

type ListProductsResponse struct {
	DataInfo        products.DataInfo      `json:"info"`
	TransactionData []TransactionsResponse `json:"data"`
}

type ListProductsResponseMerchant struct {
	DataInfo        products.DataInfo              `json:"info"`
	TransactionData []TransactionsResponseMerchant `json:"data"`
}

type MidtransConfig struct {
	MerchantId string
	ClientKey  string
	ServerKey  string
	AuthKey    string
}

type TransferBNIRequest struct {
	PaymentType            string             `json:"payment_type"`
	TransactionDetailsData TransactionDetails `json:"transaction_details"`
	BankTransferData       BankTransfer       `json:"bank_transfer"`
}

type TransactionDetails struct {
	OrderId     string `json:"order_id"`
	GrossAmount string `json:"gross_amount"`
}

type BankTransfer struct {
	Bank string `json:"bank"`
}

type VirtualAccount struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}

type PaymentAmout struct {
	PaidAt string `json:"paid_at"`
	Amount string `json:"amount"`
}
type TransferBNIResponse struct {
	StatusCode         string           `json:"status_code,omitempty"`
	StatusMessage      string           `json:"status_message,omitempty"`
	TransactionId      string           `json:"transaction_id,omitempty"`
	OrderId            string           `json:"order_id,omitempty"`
	MerchantId         string           `json:"merchant_id,omitempty"`
	GrossAmount        string           `json:"gross_amount,omitempty"`
	Currency           string           `json:"currency,omitempty"`
	PaymentType        string           `json:"payment_type,omitempty"`
	TransactionTime    string           `json:"transaction_time,omitempty"`
	TransactionStatus  string           `json:"transaction_status,omitempty"`
	VirtualAccountData []VirtualAccount `json:"va_numbers,omitempty"`
	FraudStatus        string           `json:"fraud_status,omitempty"`
	ID                 string           `json:"id,omitempty"`
}

type TransferPermataResponse struct {
	StatusCode        string `json:"status_code,omitempty"`
	StatusMessage     string `json:"status_message,omitempty"`
	TransactionId     string `json:"transaction_id,omitempty"`
	OrderId           string `json:"order_id,omitempty"`
	MerchantId        string `json:"merchant_id,omitempty"`
	GrossAmount       string `json:"gross_amount,omitempty"`
	Currency          string `json:"currency,omitempty"`
	PaymentType       string `json:"payment_type,omitempty"`
	TransactionTime   string `json:"transaction_time,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty"`
	PermataVANumber   string `json:"permata_va_number,omitempty"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	ID                string `json:"id,omitempty"`
}

type PaymentCallbackRequest struct {
	VirtualAccountData []VirtualAccount `json:"va_numbers"`
	TransactionTime    string           `json:"transaction_time"`
	TransactionStatus  string           `json:"transaction_status"`
	TransactionId      string           `json:"transaction_id"`
	StatusMessage      string           `json:"status_message"`
	StatusCode         string           `json:"status_code"`
	SignatureKey       string           `json:"signature_key"`
	SettlementTime     string           `json:"settlement_time"`
	PaymentType        string           `json:"payment_type"`
	PaymentAmoutData   []PaymentAmout   `json:"payment_amounts"`
	OrderId            string           `json:"order_id"`
	MerchantId         string           `json:"merchant_id"`
	GrossAmount        string           `json:"gross_amount"`
	FraudStatus        string           `json:"fraud_status"`
	Currency           string           `json:"currency"`
}

func (m *PaymentCallbackRequest) Bind(r *http.Request) error {
	return m.ValidatePaymentCallbackRequest()
}

func (m *PaymentCallbackRequest) ValidatePaymentCallbackRequest() error {
	return validation.ValidateStruct(m)
}

type GetPaymentInstructionRequest struct {
	InvoiceId string `json:"invoice_id"`
}

func (m *GetPaymentInstructionRequest) Bind(r *http.Request) error {
	return m.ValidateGetPaymentInstructionRequest()
}

func (m *GetPaymentInstructionRequest) ValidateGetPaymentInstructionRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.InvoiceId, validation.Required))
}

type GetPaymentInstructionResponse struct {
	PaymentMethodDesc    string   `json:"payment_method_desc"`
	PaymentMethodIcon    string   `json:"payment_method_icon"`
	TotalAmountFormatted string   `json:"total_amount_formatted"`
	VirtualAccNumber     string   `json:"virtual_account"`
	AtmTitle             string   `json:"atm_title"`
	AtmInstruction       []string `json:"instructions_atm"`
	IBInstruction        []string `json:"instructions_internet_banking"`
	MBInstuction         []string `json:"instructions_mobile_banking"`
}

type ProductsDataPageJadwal struct {
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
	JenisBan             string  `json:"jenis_ban"`
}

type GetTransactionsDetailResponse struct {
	InvoiceId          string                   `json:"invoice_id"`
	BannerInformation  string                   `json:"banner_information"`
	InstallationtTime  string                   `json:"installation_time"`
	OutletName         string                   `json:"outlet_name"`
	CsNumber           string                   `json:"outlet_cs_number"`
	OutletAddress      string                   `json:"outlet_address"`
	RescheduleTime     string                   `json:"reschedule_time"`
	IsEnableReview     bool                     `json:"is_enable_review"`
	IsEnableReschedule bool                     `json:"is_enable_reschedule"`
	PaymentMethod      string                   `json:"payment_method"`
	PaymentMethodDesc  string                   `json:"payment_method_desc"`
	PaymentMethodIcon  string                   `json:"payment_method_icon"`
	ListProduct        []ProductsDataPageJadwal `json:"list_product"`
}

type SubmitTransactionResponse struct {
	Status    string `json:"status"`
	InvoiceId string `json:"invoice_id"`
}

type GetSummaryTransactionCountResponse struct {
	WaitingPayment int `json:"waiting_payment"`
	WaitingProcess int `json:"waiting_to_process"`
	OnProgress     int `json:"on_progress"`
	Succedd        int `json:"succeed"`
}

type TransactionCommonRequest struct {
	InvoiceId string `json:"invoice_id"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
}

func (m *TransactionCommonRequest) Bind(r *http.Request) error {
	return m.ValidateTransactionCommonRequest()
}

func (m *TransactionCommonRequest) ValidateTransactionCommonRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.InvoiceId, validation.Required),
		validation.Field(&m.Status, validation.Required))
}

type GetTransactionsDetaiMerchantlResponse struct {
	InvoiceId          string                   `json:"invoice_id"`
	Status             string                   `json:"status"`
	BannerInformation  string                   `json:"banner_information"`
	InstallationtTime  string                   `json:"installation_time"`
	OutletName         string                   `json:"outlet_name"`
	CsNumber           string                   `json:"outlet_cs_number"`
	OutletAddress      string                   `json:"outlet_address"`
	RescheduleTime     string                   `json:"reschedule_time"`
	IsEnableReview     bool                     `json:"is_enable_review"`
	IsEnableReschedule bool                     `json:"is_enable_reschedule"`
	PaymentMethod      string                   `json:"payment_method"`
	PaymentMethodDesc  string                   `json:"payment_method_desc"`
	PaymentMethodIcon  string                   `json:"payment_method_icon"`
	CustomerName       string                   `json:"customer_name"`
	CustomerPhone      string                   `json:"customer_phone"`
	CustomerEmail      string                   `json:"customer_email"`
	ListProduct        []ProductsDataPageJadwal `json:"list_product"`
}

type FCMRequest struct {
	To               string       `json:"to"`
	NotificationData Notification `json:"notification"`
	FCMData          DataFCM      `json:"data"`
}

type Notification struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}

type DataFCM struct {
	Action    string `json:"action"`
	InvoiceID string `json:"invoice_id"`
}
