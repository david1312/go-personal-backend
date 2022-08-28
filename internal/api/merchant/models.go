package merchant

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (m *LoginRequest) Bind(r *http.Request) error {
	return m.ValidateLoginRequest()
}

func (m *LoginRequest) ValidateLoginRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Username, validation.Required),
		validation.Field(&m.Password, validation.Required),
	)
}

type MerchantDataResponse struct {
	Username       string `json:"username"`
	OutletId       int    `json:"outlet_id"`
	OutletName     string `json:"outlet_name"`
	OutletAvatar   string `json:"outlet_avatar"`
	OutletEmail    string `json:"outlet_email"`
	CsNumber       string `json:"outlet_cs_number"`
	OutletAddress  string `json:"outlet_address"`
	OutletCity     string `json:"outlet_city"`
	OutletGmapUrl  string `json:"outlet_gmap_url"`
}