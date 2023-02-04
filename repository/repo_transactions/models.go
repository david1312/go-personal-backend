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

type GetListTransactionsParam struct {
	Limit              int
	Page               int
	StatusTransactions []string
	CustomerId         int
}

type Transactions struct {
	InvoiceId         string    `json:"invoice_id"`
	Status            string    `json:"status"`
	TotalAmount       float64   `json:"total_amount"`
	PaymentMethod     string    `json:"payment_method"`
	PaymentMethodDesc string    `json:"payment_method_desc"`
	PaymentMethodIcon string    `json:"payment_method_icon"`
	CreatedAt         time.Time `json:"created_at"`
	PaymentDue        time.Time
	VirtualAccount    string
	OutletId          int    `json:"outlet_id"`
	OutletName        string `json:"outlet_name"`
	InstallationDate  string `json:"installation_date"`
	InstallationTime  string `json:"installation_time"`
}

type ProductsData struct {
	InvoiceId    string  `json:"invoice_id"`
	KodePLU      int32   `json:"id"`
	NamaBarang   string  `json:"nama_barang"`
	NamaUkuran   string  `json:"ukuran"`
	Qty          int     `json:"qty"`
	Harga        float64 `json:"harga"`
	HargaTotal   float64 `json:"harga_total"`
	Deskripsi    string  `json:"deskripsi"`
	DisplayImage string  `json:"display_image"`
	JenisBan     string
}

type GetTransactionsDetailData struct {
	InvoiceId         string `json:"invoice_id"`
	InstallationDate  string `json:"installation_date"`
	InstallationTime  string `json:"installation_time"`
	OutletName        string `json:"outlet_name"`
	CsNumber          string `json:"outlet_cs_number"`
	OutletAddress     string `json:"outlet_address"`
	OutletDistrict    string `json:"outlet_district"`
	OutletCity        string `json:"outlet_city"`
	Status            string `json:"status"`
	PaymentMethod     string `json:"payment_method"`
	PaymentMethodDesc string `json:"payment_method_desc"`
	PaymentMethodIcon string `json:"payment_method_icon"`
}

type GetSummaryTransactionCount struct {
	WaitingPayment int `json:"waiting_payment"`
	WaitingProcess int `json:"waiting_to_process"`
	OnProgress     int `json:"on_progress"`
	Succedd        int `json:"succeed"`
}

type FCMToken struct {
	DeviceToken string `json:"invoice_id"`
}
