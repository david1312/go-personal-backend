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

type GetAllSalesDetailRequest struct {
	NoPesanan string `json:"no_pesanan"`
}

func (m *GetAllSalesDetailRequest) Bind(r *http.Request) error {
	return m.ValidateGetAllSalesDetailRequest()
}

func (m *GetAllSalesDetailRequest) ValidateGetAllSalesDetailRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.NoPesanan, validation.Required),
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

type SalesDetailResponse struct {
	ID                  int     `json:"id"`
	NoPesanan           string  `json:"no_pesanan"`
	NoRef               string  `json:"ref"`
	Tanggal             string  `json:"tanggal"`
	NamaToko            string  `json:"nama_toko"`
	Channel             string  `json:"channel"`
	Pelanggan           string  `json:"pelanggan"`
	Status              string  `json:"status"`
	SubTotal            float64 `json:"sub_total"`
	Diskon              float64 `json:"diskon"`
	DiskonLainnya       float64 `json:"diskon_lainnya"`
	BiayaLain           float64 `json:"biaya_lain"`
	NettSales           float64 `json:"nett_sales"`
	HPP                 float64 `json:"hpp"`
	GrossProfit         float64 `json:"gross_profit"`
	PotonganMarketplace float64 `json:"potongan_marketplace"`
	NetProfit           float64 `json:"net_profit"`
}

type SalesItem struct {
	ItemId              string  `json:"item_id"`
	SKU                 string  `json:"sku"`
	NamaBarang          string  `json:"nama_barang"`
	HPPSatuan           float64 `json:"hpp_satuan"`
	SellPrice           float64 `json:"sell_price"`
	Qty                 float64 `json:"qty"`
	TotalHarga          float64 `json:"harga_satuan"`
	DiskonNumber        float64 `json:"diskon_number"`
	Diskon              float64 `json:"diskon"`
	HargaFinal          float64 `json:"harga_final"`
	HPP                 float64 `json:"hpp"`
	GrossProfit         float64 `json:"gross_profit"`
	PotonganMarketplace float64 `json:"potongan_marketplace"`
	NetProfit           float64 `json:"net_profit"`
}

type ApiResponseSalesDetail struct {
	SalesDetail SalesDetailResponse `json:"sales_detail"`
	ItemList    []SalesItem         `json:"item_list"`
}
