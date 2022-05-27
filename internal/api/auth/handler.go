package auth

import (
	"net/http"
	"semesta-ban/internal/api/customers"
	localMdl "semesta-ban/internal/api/middleware"
	"semesta-ban/internal/api/response"
	"time"
)

type AuthHandler struct {
	jwt *localMdl.JWT
}

//todo REMEMBER 30 May gmail tidak support lagi less secure app find solution

func NewAuthHandler(jwt *localMdl.JWT) *AuthHandler {
	return &AuthHandler{jwt: jwt}
}

func (usr *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	//generate token
	expiredTime := time.Now().Add(3 * time.Hour)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      authData.Uid,
		CustName: authData.CustName,
		Expired:  expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Hour * 24 * 7)
	_, tokenRefresh, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      authData.Uid,
		CustName: authData.CustName,
		Expired:  expiredTimeRefresh,
	})

	response.Yay(w, r, customers.LoginResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: tokenRefresh,
		RTExpired:    expiredTimeRefresh,
	}, http.StatusOK)

}
