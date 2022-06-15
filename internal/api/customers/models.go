package customers

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

type GetCustomerResponse struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Phone         string `json:"phone"`
	PhoneVerified bool   `json:"phone_verified"`
	Gender        string `json:"gender"`
	Avatar        string `json:"avatar"`
	Birthdate     string `json:"birthdate"`
}

type ChangePwdRequest struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

func (m *ChangePwdRequest) Bind(r *http.Request) error {
	return m.ValidateChangePwdRequest()
}

func (m *ChangePwdRequest) ValidateChangePwdRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.NewPassword, validation.Required),
		validation.Field(&m.OldPassword, validation.Required),
	)
}

type ResendEmailRequest struct {
	Email string `json:"email"`
}

func (m *ResendEmailRequest) Bind(r *http.Request) error {
	return m.ValidateResendEmailRequest()
}

func (m *ResendEmailRequest) ValidateResendEmailRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Email, validation.Required),
	)
}

type ChangeEmailRequest struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
	Code     string `json:"code"`
}

func (m *ChangeEmailRequest) Bind(r *http.Request) error {
	return m.ValidateChangeEmailRequest()
}

func (m *ChangeEmailRequest) ValidateChangeEmailRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.OldEmail, validation.Required),
		validation.Field(&m.NewEmail, validation.Required),
		validation.Field(&m.Code, validation.Required),
	)
}

type UpdateNameRequest struct {
	Name string `json:"name"`
}

func (m *UpdateNameRequest) Bind(r *http.Request) error {
	return m.ValidateUpdateNameRequest()
}

func (m *UpdateNameRequest) ValidateUpdateNameRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Name, validation.Required),
	)
}

type UpdateGenderRequest struct {
	Gender string `json:"gender"`
}

func (m *UpdateGenderRequest) Bind(r *http.Request) error {
	return m.ValidateUpdateGenderRequest()
}

func (m *UpdateGenderRequest) ValidateUpdateGenderRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Gender, validation.Required),
	)
}

type UpdatePhoneRequest struct {
	Phone string `json:"phone"`
}

func (m *UpdatePhoneRequest) Bind(r *http.Request) error {
	return m.ValidateUpdatePhoneRequest()
}

func (m *UpdatePhoneRequest) ValidateUpdatePhoneRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Phone, validation.Required),
	)
}

type UpdateBirthDateRequest struct {
	Birthdate string `json:"birthdate"`
}

func (m *UpdateBirthDateRequest) Bind(r *http.Request) error {
	return m.ValidateUpdateBirthDateRequest()
}

func (m *UpdateBirthDateRequest) ValidateUpdateBirthDateRequest() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Birthdate, validation.Required),
	)
}
