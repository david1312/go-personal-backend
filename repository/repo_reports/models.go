package repo_reports

type Sales struct {
	ID                  int
	Tanggal             string
	TipeTransaksi       string
	Ref                 string
	NoPesanan           string
	Status              string
	Channel             string
	NamaToko            string
	Pelanggan           string
	SubTotal            float64
	Diskon              float64
	DiskonLainnya       float64
	PotonganBiaya       float64
	BiayaLain           float64
	TermasukPajak       string
	Pajak               float64
	Ongkir              float64
	Asuransi            float64
	NettSales           float64
	HPP                 float64
	GrossProfit         float64
	PotonganMarketPlace float64
	NetProfit           float64
}
