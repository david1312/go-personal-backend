package users

import (
	"errors"
	"fmt"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/pkg/log"
	custRepo "semesta-ban/repository/repo_customers"
	"time"

	localMdl "semesta-ban/internal/api/middleware"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"gopkg.in/gomail.v2"
)

//todo move to config
const CONFIG_SMTP_HOST = "mail.sunmorisemestaban.com"
const CONFIG_SMTP_PORT = 465
const CONFIG_SENDER_NAME = "PT. Sunmori Semesta Ban <support@sunmorisemestaban.com>"
const CONFIG_AUTH_EMAIL = "support@sunmorisemestaban.com"
const CONFIG_AUTH_PASSWORD = "spyxfamily13"
const CONFIG_API_URL = "https://api.sunmorisemestaban.com"

type UsersHandler struct {
	db             *sqlx.DB
	custRepository custRepo.CustomersRepository
	jwt            *localMdl.JWT
}

//todo REMEMBER 30 May gmail tidak support lagi less secure app find solution

func NewUsersHandler(db *sqlx.DB, cr custRepo.CustomersRepository, jwt *localMdl.JWT) *UsersHandler {
	return &UsersHandler{db: db, custRepository: cr, jwt: jwt}
}

func (usr *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		p   RegisterRequest
		ctx = r.Context()
	)
	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	isEmailPhoneExist, errCode, err := usr.custRepository.CheckEmailExist(ctx, p.Email)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if isEmailPhoneExist {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrEmailExists), crashy.ErrEmailExists, crashy.Message(crashy.ErrEmailExists)), http.StatusBadRequest)
		return
	}

	hashedTokenEmail := helper.GenerateHashString()
	errCode, err = usr.custRepository.Register(ctx, p.Name, p.Email, hashedTokenEmail, p.Password)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	bodyEmail := "Hallo <b>" + p.Name + "</b>!, <br> Terimakasih telah bersedia bergabung bersama kami, silahkan lakukan verifikasi email anda dengan klik link berikut : " + CONFIG_API_URL + "/v1/verify?val=" + hashedTokenEmail
	_ = sendMail(p.Email, "Selamat Menjadi Bagian Pengguna Semesta Ban!", bodyEmail) // keep going even though send email failed

	//generate token
	expiredTime := time.Now().Add(60 * time.Minute)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		CustName: p.Name,
		Expired:  expiredTime,
	})

	response.Yay(w, r, LoginResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: "//todo implement refresh token",
		RTExpired:    expiredTime,
	}, http.StatusOK)

}

func (usr *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		p   LoginRequest
		ctx = r.Context()
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	name, errCode, err := usr.custRepository.Login(ctx, p.Email, p.Password)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//generate token
	expiredTime := time.Now().Add(60 * time.Minute)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		CustName: name,
		Expired:  expiredTime,
	})

	response.Yay(w, r, LoginResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: "//todo implement refresh token",
		RTExpired:    expiredTime,
	}, http.StatusOK)
	// response.Nay(w, r, crashy.New(errors.New("real error"), crashy.ErrCodeFormatting, "output error"), http.StatusBadRequest)
}

func (usr *UsersHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	response.Yay(w, r, "profile", http.StatusOK)
}

func (usr *UsersHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var (
		fp  = NewVerifyEmailParam(r)
		ctx = r.Context()
	)
	if len(fp.TokenEmail) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCode(crashy.ErrCodeValidation), crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}

	errCode, err := usr.custRepository.VerifyEmail(ctx, fp.TokenEmail)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Selamat! Email anda berhasil diverifikasi. Silahkan buka aplikasi semesta ban untuk melanjutkan.")
}

func sendMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", CONFIG_SENDER_NAME)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send the email to Bob
	d := gomail.NewDialer(CONFIG_SMTP_HOST, CONFIG_SMTP_PORT, CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD)
	if err := d.DialAndSend(m); err != nil {
		log.Error(err)
		return err
	}
	log.Infof("success sent email to %v", to)

	return nil
}
