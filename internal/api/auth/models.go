package auth

import "time"

type AnonymousToken struct {
	AnonToken string    `json:"anonymous_token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type RefreshTokenResponse struct {
	Token        string    `json:"token"`
	ExpiredAt    time.Time `json:"expired_at"`
	RefreshToken string    `json:"refresh_token"`
	RTExpired    time.Time `json:"refresh_token_expired"`
	AnonToken    string    `json:"anonymous_token"`
	AnonExpired  time.Time `json:"anonymous_token_expired"`
}

type Version struct {
	ApiVersion string `json:"api_version"`
}
