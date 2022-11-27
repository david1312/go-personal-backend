package repo_products

type Products struct {
	KodePLU            int32   `json:"kode_plu"`
	KodeBarang         string  `json:"kode_barang"`
	NamaBarang         string  `json:"nama_barang"`
	Barcode            string  `json:"barcode"`
	UnitJual           string  `json:"unit_jual"`
	Qty                int     `json:"qty"`
	KodeSupplier       string  `json:"kode_supplier"`
	NamaSupplier       string  `json:"nama_supplier"`
	HargaJual          float64 `json:"harga_jual"`
	HargaJualFinal     float64 `json:"harga_jual_final"`
	Disc               float32 `json:"disc"`
	IDUkuran           string  `json:"id_ukuran"`
	IDPosisi           string  `json:"id_posisi"`
	IDMerk             string  `json:"id_merk"`
	NamaUkuran         string  `json:"nama_ukuran"`
	NamaPosisi         string  `json:"nama_posisi"`
	NamaMerk           string  `json:"nama_merk"`
	JenisBan           string  `json:"jenis_ban"`
	DisplayImage       string  `json:"display_image"`
	JenisMotor         string  `json:"jenis_motor"`
	TotalTerjual       int     `json:"total_terjual"`
	Deskripsi          string  `json:"deskripsi"`
	IsWishlist         bool    `json:"is_wishlist"`
	CartItemId         int     `json:"cart_item_id"`
	CartItemQty        int     `json:"cart_item_qty"`
	CartItemIsSelected bool    `json:"cart_item_is_selected"`
	StockAll           int     `json:"stock_all"`
	IDRingBan          int
}

type ProductsParamsTemp struct {
	Limit      int
	Page       int
	Name       string
	UkuranBan  []string
	Posisi     string
	MerkBan    []string
	MerkMotor  int
	IdMotor    int
	ArrUkuran  []string
	MinPrice   int
	MaxPrice   int
	OrderBy    string
	OrderType  string
	CustomerId int
}

type ProductImage struct {
	Id             int `json:"id"`
	Url            string
	IsDisplay      string
	Count          int
	KodeBarang     string
	IsDisplayFixed bool
}

type CustomerResponse struct {
	Id  int    `json:"id"`
	Uid string `json:"uid"`
}

type ListUkuranBan struct {
	Id int `json:"id"`
}

type MotorCycleCompatibility struct {
	Id           int
	Motor        string
	DisplayImage string
}

type ProductReview struct {
	Name    string `json:"url"`
	Avatar  string `json:"avatar"`
	Date    string `json:"is_display"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type UpdateProductParam struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IdTIreBrand string `json:"id_tire_brand"`
	TireType    string `json:"tire_type"`
	Size        string `json:"size"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
}
