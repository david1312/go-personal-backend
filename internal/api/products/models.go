package products

import (
	"net/http"
	"strconv"
)

type ProductsParams struct {
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

type DataInfo struct {
	CurrentPage int `json:"cur_page"`
	MaxPage     int `json:"max_page"`
	Limit       int `json:"limit"`
	TotalRecord int `json:"total_record"`
}

func NewProductsParams(r *http.Request) ProductsParams {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if page < 1 {
		page = 1
	}

	minPrice, _ := strconv.Atoi(r.URL.Query().Get("minprice"))
	maxPrice, _ := strconv.Atoi(r.URL.Query().Get("maxprice"))

	// add filetype validation

	return ProductsParams{
		Limit:     limit,
		Page:      page,
		Name:      r.URL.Query().Get("name"),
		UkuranBan: r.URL.Query().Get("ukuran"),
		Posisi:    r.URL.Query().Get("posisi"),
		MerkBan:   r.URL.Query().Get("merkban"),
		MerkMotor: r.URL.Query().Get("merkmotor"),
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
		OrderBy:   r.URL.Query().Get("orderby"),
		OrderType: r.URL.Query().Get("ordertype"),
	}
}

type ListProductsResponse struct {
	DataInfo DataInfo           `json:"info"`
	Products []ProductsResponse `json:"data"`
}

type ProductsResponse struct {
	KodePLU        int32   `json:"id"`
	NamaBarang     string  `json:"nama_barang"`
	Disc           float32 `json:"disc"`
	NamaUkuran     string  `json:"ukuran"`
	HargaJualFinal float64 `json:"harga_jual_final"`
	IsWishList     bool    `json:"is_wishlist"`
	JenisBan       string  `json:"jenis_ban"`
	DisplayImage   string  `json:"display_image"`
}

type ProductDetailResponse struct {
	KodePLU        int32   `json:"id"`
	NamaBarang     string  `json:"nama_barang"`
	Disc           float32 `json:"disc"`
	NamaUkuran     string  `json:"ukuran"`
	HargaJualFinal float64 `json:"harga_jual_final"`
	IsWishList     bool    `json:"is_wishlist"`
	JenisBan       string  `json:"jenis_ban"`
	Posisi       string  `json:"posisi"`
	JenisMotor       string  `json:"jenis_motor"`
	TotalTerjual int			`json:"total_terjual"`
	Deskripsi string  `json:"deskripsi"`
	ImageList   []ProductImage  `json:"image_list"`
	ReviewList []ProductReview `json:"product_review"`	
	Kompatibilitas []MotorCycleCompatibility `json:"kompatibilitas"`
}

type ProductImage struct{
	Url   string  `json:"url"`
	IsDisplay   string  `json:"is_display"`
}

type ProductReview struct{
	Name   string  `json:"url"`
	Avatar string   `json:"avatar"`
	Date   string  `json:"is_display"`
	Rating  int `json:"rating"`
	Comment string `json:"comment"`
}

type MotorCycleCompatibility struct{
	MerkMotor string  `json:"url"`
	DisplayImage   string  `json:"display_image"`
}