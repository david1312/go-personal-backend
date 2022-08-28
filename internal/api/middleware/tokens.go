package middleware

import (
	"errors"
	"time"
)

type Token struct {
	Uid      string    `json:"uid"`
	CustName string    `json:"cust_name"`
	Expired  time.Time `json:"expired"`
}

type RefreshToken struct {
	Uid      string    `json:"uid"`
	CustName string    `json:"cust_name"`
	Expired  time.Time `json:"expired"`
}

type MerchantToken struct {
	OutletId int       `json:"outlet_id"`
	Username string    `json:"username"`
	Expired  time.Time `json:"expired"`
}

func (t *Token) Valid() error {
	// Cek Expired Token
	if t.Expired.Before(time.Now()) {
		return errors.New("Token Expired")
	}
	return nil
}

func (t *MerchantToken) Valid() error {
	// Cek Expired Token
	if t.Expired.Before(time.Now()) {
		return errors.New("Token Expired")
	}
	return nil
}
