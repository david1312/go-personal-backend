package auth

import (
	"net/http"
	localMdl "semesta-ban/internal/api/middleware"
	"semesta-ban/internal/api/response"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	jwt       *localMdl.JWT
	anonToken *localMdl.JWT
}

//todo REMEMBER 30 May gmail tidak support lagi less secure app find solution

func NewAuthHandler(jwt *localMdl.JWT, an *localMdl.JWT) *AuthHandler {
	return &AuthHandler{jwt: jwt, anonToken: an}
}

func (usr *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	//generate token
	expiredTime := time.Now().Add(24 * time.Hour)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      authData.Uid,
		CustName: authData.CustName,
		Expired:  expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Hour * 24 * 31)
	_, tokenRefresh, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      authData.Uid,
		CustName: authData.CustName,
		Expired:  expiredTimeRefresh,
	})

	expiredTimeAnon := time.Now().Add(time.Hour * 24 * 30)
	_, anonToken, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      uuid.New().String(),
		CustName: "",
		Expired:  expiredTime,
	})


	response.Yay(w, r, RefreshTokenResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: tokenRefresh,
		RTExpired:    expiredTimeRefresh,
		AnonToken: anonToken,
		AnonExpired: expiredTimeAnon,
	}, http.StatusOK)

}

func (usr *AuthHandler) GetAnonymousToken(w http.ResponseWriter, r *http.Request) {
	// var (
	// 	ctx      = r.Context()
	// 	// authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	// )

	//generate token
	expiredTime := time.Now().Add(time.Hour * 24 * 30)
	_, token, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      uuid.New().String(),
		CustName: "",
		Expired:  expiredTime,
	})

	response.Yay(w, r, AnonymousToken{
		AnonToken: token,
		ExpiredAt: expiredTime,
	}, http.StatusOK)

}
