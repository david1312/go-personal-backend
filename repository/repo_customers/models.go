package repo_customers

import (
	"database/sql"
)

type Customers struct {
	ID                 int32          `json:"id"`
	CustId             string         `json:"cust_id"`
	Uid                string         `json:"uid"`
	Name               string         `json:"name"`
	Password           string         `json:"password"`
	Email              string         `json:"email"`
	EmailVerifiedToken sql.NullString `json:"email_verified_token"`
	EmailVerifiedAt    sql.NullTime   `json:"email_verified_at"`
	Gender             sql.NullString `json:"gender"`
	IsActive           bool           `json:"is_active"`
	Phone              sql.NullString `json:"phone"`
	PhoneVerifiedAt    sql.NullTime   `json:"phone_verified_at"`
	Avatar             sql.NullString `json:"avatar"`
	Birthdate          sql.NullTime   `json:"birthdate"`
}

type InsertCustomerParam struct {
	Name               string `json:"name"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	EmailVerifiedToken string `json:"email_verified_token"`
}
