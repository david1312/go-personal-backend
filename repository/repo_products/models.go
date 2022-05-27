package repo_products

type Products struct {
	KodePLU        int32   `json:"kode_plu"`
	KodeBarang     string  `json:"kode_barang"`
	NamaBarang     string  `json:"nama_barang"`
	Barcode        string  `json:"barcode"`
	UnitJual       string  `json:"unit_jual"`
	Qty            float32 `json:"qty"`
	KodeSupplier   string  `json:"kode_supplier"`
	NamaSupplier   string  `json:"nama_supplier"`
	HargaJual      float64 `json:"harga_jual_final"`
	HargaJualFinal float64 `json:"harga_jual"`
	Disc           float32 `json:"disc"`
	IDUkuran       string  `json:"id_ukuran"`
	IDPosisi       string  `json:"id_posisi"`
	IDMerk         string  `json:"id_merk"`
	NamaUkuran     string  `json:"nama_ukuran"`
	NamaPosisi     string  `json:"nama_posisi"`
	NamaMerk       string  `json:"nama_merk"`
}

type ProductsParamsTemp struct {
	Limit     int
	Page      int
	Name      string
	UkuranBan string
	Posisi    string
	MerkBan   string
	MerkMotor string
	MinPrice  int
	MaxPrice  int
	OrderBy   string
	OrderType string
}
