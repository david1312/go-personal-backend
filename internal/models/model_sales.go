package models

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GetAllSalesRequest struct {
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	NoPesanan string `json:"no_pesanan"`
}

func (m *GetAllSalesRequest) Bind(r *http.Request) error {
	return m.ValidateGetAllSalesRequest()
}

func (m *GetAllSalesRequest) ValidateGetAllSalesRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.StartDate, validation.Required),
		validation.Field(&m.EndDate, validation.Required),
	)
}

type Pagination struct {
	CurrentPage int `json:"cur_page"`
	MaxPage     int `json:"max_page"`
	Limit       int `json:"limit"`
	TotalRecord int `json:"total_record"`
}

type SummarySales struct {
	TotalNettSales           float64 `json:"total_nett_sales"`
	TotalGross               float64 `json:"total_gross_profit"`
	TotalPotonganMarketplace float64 `json:"total_potongan_marketplace"`
	TotalNetProfit           float64 `json:"total_net_profit"`
}

type SalesResponse struct {
	ID                  int     `json:"id"`
	Tanggal             string  `json:"tanggal"`
	NoPesanan           string  `json:"no_pesanan"`
	Status              string  `json:"status"`
	Channel             string  `json:"channel"`
	NettSales           float64 `json:"nett_sales"`
	GrossProfit         float64 `json:"gross_profit"`
	PotonganMarketplace float64 `json:"potongan_marketplace"`
	NetProfit           float64 `json:"net_profit"`
}

type ApiResponseSales struct {
	PaginationData Pagination      `json:"pagination"`
	SummaryData    SummarySales    `json:"summary_data"`
	SalesList      []SalesResponse `json:"data"`
}
