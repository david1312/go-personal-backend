package products

import (
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GetProductsRequest struct {
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	Name      string `json:"name"`
	UkuranBan []string `json:"ukuran"`
	MerkBan   []string `json:"merkban"`
	MerkMotor []int  `json:"merkmotor"`
	IdMotor   int    `json:"idmotor"`
	MinPrice  int    `json:"minprice"`
	MaxPrice  int    `json:"maxprice"`
	OrderBy   string `json:"orderby"`
}

func (m *GetProductsRequest) Bind(r *http.Request) error {
	return m.ValidateGetProductsRequest()
}

func (m *GetProductsRequest) ValidateGetProductsRequest() error {
	return validation.ValidateStruct(m)
}

type ProductsParams struct {
	Limit     int
	Page      int
	Name      string
	UkuranBan string
	Posisi    string
	MerkBan   string
	MerkMotor int
	IdMotor   int
	MinPrice  int
	MaxPrice  int
	OrderBy   string
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

	merkMotor, _ := strconv.Atoi(r.URL.Query().Get("merkmotor"))
	idMotor, _ := strconv.Atoi(r.URL.Query().Get("idmotor"))

	// add filetype validation

	return ProductsParams{
		Limit:     limit,
		Page:      page,
		Name:      r.URL.Query().Get("name"),
		UkuranBan: r.URL.Query().Get("ukuran"),
		Posisi:    r.URL.Query().Get("posisi"),
		MerkBan:   r.URL.Query().Get("merkban"),
		MerkMotor: merkMotor,
		IdMotor:   idMotor,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
		OrderBy:   r.URL.Query().Get("orderby"),
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
	KodePLU        int32                     `json:"id"`
	NamaBarang     string                    `json:"nama_barang"`
	Disc           float32                   `json:"disc"`
	NamaUkuran     string                    `json:"ukuran"`
	HargaJualFinal float64                   `json:"harga_jual_final"`
	IsWishList     bool                      `json:"is_wishlist"`
	JenisBan       string                    `json:"jenis_ban"`
	Posisi         string                    `json:"posisi"`
	JenisMotor     string                    `json:"jenis_motor"`
	TotalTerjual   int                       `json:"total_terjual"`
	Deskripsi      string                    `json:"deskripsi"`
	ImageList      []ProductImage            `json:"image_list"`
	ReviewList     []ProductReview           `json:"product_review"`
	Kompatibilitas []MotorCycleCompatibility `json:"kompatibilitas"`
}

type ProductImage struct {
	Url       string `json:"url"`
	IsDisplay string `json:"is_display"`
}

type ProductReview struct {
	Name    string `json:"url"`
	Avatar  string `json:"avatar"`
	Date    string `json:"is_display"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type MotorCycleCompatibility struct {
	MerkMotor    string `json:"url"`
	DisplayImage string `json:"display_image"`
}

type WishlistRequest struct {
	KodePLU int `json:"id"`
}

func (m *WishlistRequest) Bind(r *http.Request) error {
	return m.ValidateWishlistRequest()
}

func (m *WishlistRequest) ValidateWishlistRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.KodePLU, validation.Required),
	)
}

type CartItemUpdateRequest struct {
	CartItemId int  `json:"cart_item_id"`
	Qty        int  `json:"qty"`
	IsSelected bool `json:"is_selected"`
}

func (m *CartItemUpdateRequest) Bind(r *http.Request) error {
	return m.ValidateCartItemUpdateRequest()
}

func (m *CartItemUpdateRequest) ValidateCartItemUpdateRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.CartItemId, validation.Required),
		validation.Field(&m.Qty, validation.Required),
	)
}

type CartItemRemoveRequest struct {
	CartItemId int `json:"cart_item_id"`
}

func (m *CartItemRemoveRequest) Bind(r *http.Request) error {
	return m.ValidateCartItemRemoveRequest()
}

func (m *CartItemRemoveRequest) ValidateCartItemRemoveRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.CartItemId, validation.Required),
	)
}

type CartSelectAllRequest struct {
	CartId      int  `json:"cart_id"`
	IsSelectAll bool `json:"is_select_all"`
}

func (m *CartSelectAllRequest) Bind(r *http.Request) error {
	return m.ValidateCartSelectAllRequest()
}

func (m *CartSelectAllRequest) ValidateCartSelectAllRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.CartId, validation.Required),
	)
}

type CartResponse struct {
	CartItemId         int     `json:"cart_item_id"`
	CartItemQty        int     `json:"cart_item_qty"`
	CartItemIsSelected bool    `json:"cart_item_is_selected"`
	KodePLU            int32   `json:"id"`
	NamaBarang         string  `json:"nama_barang"`
	Disc               float32 `json:"disc"`
	NamaUkuran         string  `json:"ukuran"`
	HargaJualFinal     float64 `json:"harga_jual_final"`
	JenisBan           string  `json:"jenis_ban"`
	DisplayImage       string  `json:"display_image"`
}

type ListItemCartResponse struct {
	DataInfo  DataInfo       `json:"info"`
	CartId    int            `json:"cart_id"`
	CartsItem []CartResponse `json:"data"`
}
