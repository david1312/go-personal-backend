package models

type JULoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"userName"`
}

type JULoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
