package users

import (
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *LoginRequest) Bind(r *http.Request) error {
	return m.ValidateLoginRequest()
}

func (m *LoginRequest) ValidateLoginRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Email, validation.Required),
		validation.Field(&m.Password, validation.Required),
	)
}

type LoginResponse struct {
	Token        string    `json:"token"`
	ExpiredAt    time.Time `json:"expired_at"`
	RefreshToken string    `json:"refresh_token"`
	RTExpired    time.Time `json:"refresh_token_expired"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *RegisterRequest) Bind(r *http.Request) error {
	return m.ValidateRegisterRequest()
}

func (m *RegisterRequest) ValidateRegisterRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Email, validation.Required),
		validation.Field(&m.Password, validation.Required),
	)
}

type VerifyEmailParam struct {
	TokenEmail string
}

func NewVerifyEmailParam(r *http.Request) VerifyEmailParam {
	return VerifyEmailParam{
		TokenEmail: r.URL.Query().Get("val"),
	}
}
