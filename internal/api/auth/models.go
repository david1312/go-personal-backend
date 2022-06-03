package auth

import "time"

type AnonymousToken struct {
	AnonToken string    `json:"anonymous_token"`
	ExpiredAt time.Time `json:"expired_at"`
}
