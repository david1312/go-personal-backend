package middleware

import (
	"errors"
	"time"
)

type Token struct {
	CustName string    `json:"cust_name"`
	Expired  time.Time `json:"expired"`
}

func (t *Token) Valid() error {
	// Cek Expired Token
	if t.Expired.Before(time.Now()) {
		return errors.New("Token Expired")
	}
	return nil
}
