package repo_customers

import (
	"database/sql"
)

type Customers struct {
	ID                 int32          `json:"id"`
	Name               string         `json:"name"`
	Password           string         `json:"password"`
	Email              int32          `json:"email"`
	EmailVerifiedToken sql.NullString `json:"email_verified_token"`
	EmailVerifiedAt    sql.NullTime   `json:"email_verified_at"`
	Gender             string         `json:"gender"`
	IsActive           bool           `json:"is_active"`
	Phone              string         `json:"phone"`
	PhoneVerifiedAt    sql.NullTime   `json:"phone_verified_at"`
	Avatar             sql.NullString `json:"avatar"`
}

type InsertCustomerParam struct {
	Name               string `json:"name"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	EmailVerifiedToken string `json:"email_verified_token"`
}
