package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/crashy"
	"strings"

	"github.com/go-chi/jwtauth"
)

type ContextKey string

type GuardType int32

const (
	GuardAnonymous GuardType = iota
	GuardAccess

	AuthPrefix = "Bearer "
	CtxKey     = ContextKey("context-data")
)

type JWT struct {
	*jwtauth.JWTAuth
}

func New(secret []byte) *JWT {
	return &JWT{jwtauth.New("HS256", secret, nil)}
}

func (j *JWT) Encode(token Token) string {
	_, tokenString, _ := j.JWTAuth.Encode(&token)
	return tokenString
}

func (j *JWT) Decode(tokenString string) (token Token, err error) {
	jwtToken, err := j.JWTAuth.Decode(tokenString)
	fmt.Println(jwtToken.Claims)
	if err != nil {
		return
	}
	return
}

func (j *JWT) GetToken(r *http.Request) Token {
	token, ok := r.Context().Value("token").(Token)
	if !ok {
		return Token{}
	}
	return token
}

func (j *JWT) AuthMiddleware(typ GuardType) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			authHeader := r.Header.Get("Authorization")

			if len(authHeader) == 0 || !strings.HasPrefix(authHeader, AuthPrefix) {
				response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCodeUnauthorized, "invalid token"), http.StatusUnauthorized)
				return
			}

			jwtToken, err := jwtauth.VerifyRequest(j.JWTAuth, r, TokenFromHeader)
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "error verifying token"), http.StatusUnauthorized)
				return
			}

			var claims Token
			b, err := json.Marshal(jwtToken.Claims) //Encode Token
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "invalid token"), http.StatusUnauthorized)
				return
			}

			err = json.Unmarshal(b, &claims)
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeDataRead, "unable to parse token"), http.StatusBadRequest)
				return
			}

			if typ == GuardAccess && len(claims.CustName) < 1 {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "invalid token"), http.StatusUnauthorized)
				return
			}

			err = claims.Valid()
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "token expired"), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), CtxKey, claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (j *JWT) AuthMiddlewareMerchant(typ GuardType) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			authHeader := r.Header.Get("Authorization")

			if len(authHeader) == 0 || !strings.HasPrefix(authHeader, AuthPrefix) {
				response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidToken), crashy.ErrCodeUnauthorized, "invalid token"), http.StatusUnauthorized)
				return
			}

			jwtToken, err := jwtauth.VerifyRequest(j.JWTAuth, r, TokenFromHeader)
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "error verifying token"), http.StatusUnauthorized)
				return
			}

			var claims MerchantToken
			b, err := json.Marshal(jwtToken.Claims) //Encode Token
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "invalid token"), http.StatusUnauthorized)
				return
			}

			err = json.Unmarshal(b, &claims)
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeDataRead, "unable to parse token"), http.StatusBadRequest)
				return
			}

			if typ == GuardAccess && len(claims.Username) < 1 {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "invalid token"), http.StatusUnauthorized)
				return
			}

			err = claims.Valid()
			if err != nil {
				response.Nay(w, r, crashy.New(err, crashy.ErrCodeUnauthorized, "token expired"), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), CtxKey, claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TokenFromHeader(r *http.Request) string {
	return r.Header.Get("Authorization")[7:]
}
