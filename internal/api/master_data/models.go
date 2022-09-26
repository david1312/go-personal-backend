package master_data

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Outlet struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	MapUrl    string  `json:"gmap_url"`
}

type Gender struct {
	Value string `json:"value"`
}

type MerkBan struct {
	IdMerk string `json:"id_merk"`
	Merk   string `json:"merk"`
	Icon   string `json:"icon"`
}

type SortBy struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// todo update list product sort by
var ListSortBy = []SortBy{{
	Value: "latest",
	Label: "Terbaru",
},
	{
		Value: "max_price",
		Label: "Harga Tertinggi",
	},
	{
		Value: "min_price",
		Label: "Harga Terendah",
	},
}

type UkuranBan struct {
	Ukuran string `json:"ukuran"`
}

type UkuranBanTemp struct {
	RingBan string `json:"ring_ban"`
	Ukuran  string `json:"ukuran"`
}

type ListUkuranBan struct {
	RingBan    string      `json:"ring_ban"`
	ListUkuran []UkuranBan `json:"list_ukuran"`
}

type MerkMotor struct {
	Id   int    `json:"id"`
	Nama string `json:"nama"`
	Icon string `json:"icon"`
}

type Motor struct {
	Id   int    `json:"id"`
	Nama string `json:"nama"`
	Icon string `json:"icon"`
}

type ListMotor struct {
	Category  string  `json:"kategori"`
	ListMotor []Motor `json:"list_motor"`
}

type PaymentMethod struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
	Icon        string `json:"icon"`
}

type ListPaymentMethod struct {
	Category          string          `json:"category"`
	ListPaymentMethod []PaymentMethod `json:"list_payment_method"`
}

type PromoBanner struct {
	Alt      string `json:"alt"`
	ImageUrl string `json:"img_url"`
}

type ImageAssetResponse struct {
	PromoBannerData PromoBanner `json:"promo_banner"`
}

type TireType struct {
	Value string `json:"value"`
}

type MasterDataCommonRequest struct {
	Id int `json:"id"`
}

func (m *MasterDataCommonRequest) Bind(r *http.Request) error {
	return m.ValidateMasterDataCommonRequest()
}

func (m *MasterDataCommonRequest) ValidateMasterDataCommonRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Id, validation.Required),
	)
}

type MasterDataCommonRequestSec struct {
	Id string `json:"id"`
}

func (m *MasterDataCommonRequestSec) Bind(r *http.Request) error {
	return m.ValidateMasterDataCommonRequestSec()
}

func (m *MasterDataCommonRequestSec) ValidateMasterDataCommonRequestSec() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Id, validation.Required),
	)
}

type UpdateBrandMotorReq struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (m *UpdateBrandMotorReq) Bind(r *http.Request) error {
	return m.ValidateUpdateBrandMotorReq()
}

func (m *UpdateBrandMotorReq) ValidateUpdateBrandMotorReq() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Id, validation.Required),
		validation.Field(&m.Name, validation.Required),
	)
}

type UpdateTireBrandReq struct {
	Id   string    `json:"id"`
	Name string `json:"name"`
	Ranking int `json:"ranking"`
}

func (m *UpdateTireBrandReq) Bind(r *http.Request) error {
	return m.ValidateUpdateTireBrandReq()
}

func (m *UpdateTireBrandReq) ValidateUpdateTireBrandReq() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Id, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Ranking, validation.Required),
	)
}
