package auth

import (
	localMdl "libra-internal/internal/api/middleware"
	"libra-internal/internal/api/response"
	"libra-internal/pkg/constants"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	jwt       *localMdl.JWT
	anonToken *localMdl.JWT
}

func NewAuthHandler(jwt *localMdl.JWT, an *localMdl.JWT) *AuthHandler {
	return &AuthHandler{jwt: jwt, anonToken: an}
}

func (usr *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	//generate token
	expiredTime := time.Now().Add(constants.LoginTokenExpiry)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      authData.Uid,
		CustName: authData.CustName,
		Expired:  expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(constants.RefreshTokenExpiry)
	_, tokenRefresh, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      authData.Uid,
		CustName: authData.CustName,
		Expired:  expiredTimeRefresh,
	})

	//renew anon token expiration
	expiredTimeAnon := time.Now().Add(constants.AnonymousTokenExpiry)
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
		AnonToken:    anonToken,
		AnonExpired:  expiredTimeAnon,
	}, http.StatusOK)

}

func (usr *AuthHandler) GetAnonymousToken(w http.ResponseWriter, r *http.Request) {
	uid := uuid.New().String()

	//generate token
	expiredTime := time.Now().Add(constants.AnonymousTokenExpiry)
	_, token, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      uid,
		CustName: "",
		Expired:  expiredTime,
	})

	response.Yay(w, r, AnonymousToken{
		AnonToken: token,
		ExpiredAt: expiredTime,
	}, http.StatusOK)

}

func (usr *AuthHandler) GetVersion(w http.ResponseWriter, r *http.Request) {

	response.Yay(w, r, Version{
		ApiVersion: constants.ApiVersion,
	}, http.StatusOK)

}
