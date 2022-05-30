package customers

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
	cn "semesta-ban/pkg/constants"

	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	baseAssetUrl   string
}

//todo REMEMBER 30 May gmail tidak support lagi less secure app find solution

func NewUsersHandler(db *sqlx.DB, cr custRepo.CustomersRepository, jwt *localMdl.JWT, baseAssetUrl string) *UsersHandler {
	return &UsersHandler{db: db, custRepository: cr, jwt: jwt, baseAssetUrl: baseAssetUrl}
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

	if len(p.Password) < 6 { // todo implement number from config
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCode(crashy.ErrCodeValidation), crashy.Message(crashy.ErrShortPassword)), http.StatusBadRequest)
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
	uid := uuid.New().String()
	errCode, err = usr.custRepository.Register(ctx, p.Name, p.Email, hashedTokenEmail, p.Password, uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	//temporary
	bodyEmail := "Hallo <b>" + p.Name + "</b>!, <br> Terimakasih telah bersedia bergabung bersama kami, silahkan lakukan verifikasi email anda dengan klik link berikut : " + CONFIG_API_URL + "/v1/verify?val=" + hashedTokenEmail
	_ = sendMail(p.Email, "Selamat Menjadi Bagian Pengguna Semesta Ban!", bodyEmail) // keep going even though send email failed

	//generate token
	expiredTime := time.Now().Add(3 * time.Hour)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      uid,
		CustName: p.Name,

		Expired: expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Hour * 24 * 7)
	_, tokenRefresh, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      uid,
		CustName: p.Name,
		Expired:  expiredTimeRefresh,
	})

	response.Yay(w, r, LoginResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: tokenRefresh,
		RTExpired:    expiredTimeRefresh,
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

	customer, errCode, err := usr.custRepository.Login(ctx, p.Email, p.Password)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//generate token
	expiredTime := time.Now().Add(3 * time.Hour)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      customer.Uid,
		CustName: customer.Name,
		Expired:  expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Hour * 24 * 7)
	_, tokenRefresh, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      customer.Uid,
		CustName: customer.Name,
		Expired:  expiredTimeRefresh,
	})

	response.Yay(w, r, LoginResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: tokenRefresh,
		RTExpired:    expiredTimeRefresh,
	}, http.StatusOK)
}

func (usr *UsersHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)
	customer, errCode, err := usr.custRepository.GetCustomer(ctx, authData.Uid)

	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	isEmailVerified := true
	if customer.EmailVerifiedAt.Time.IsZero() {
		isEmailVerified = false
	}
	isPhoneVerified := true
	if customer.PhoneVerifiedAt.Time.IsZero() {
		isPhoneVerified = false
	}
	birthdateVal := customer.Birthdate.Time.Format("2006-01-02")
	if customer.Birthdate.Time.IsZero() {
		birthdateVal = ""
	}

	response.Yay(w, r, GetCustomerResponse{
		Name:          customer.Name,
		Email:         customer.Email,
		EmailVerified: isEmailVerified,
		Phone:         customer.Phone.String,
		PhoneVerified: isPhoneVerified,
		Gender:        customer.Gender.String,
		Avatar:        usr.baseAssetUrl + cn.UserDir + customer.Avatar.String,
		Birthdate:     birthdateVal,
	}, http.StatusOK)
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

func (usr *UsersHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var (
		p        ChangePwdRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	if len(p.NewPassword) < 6 { // todo implement number from config
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCode(crashy.ErrCodeValidation), crashy.Message(crashy.ErrShortPassword)), http.StatusBadRequest)
		return
	}

	errCode, err := usr.custRepository.ChangePassword(ctx, authData.Uid, p.OldPassword, p.NewPassword)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (usr *UsersHandler) ResendEmailVerification(w http.ResponseWriter, r *http.Request) {
	var (
		p        ResendEmailRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	emailToken, errCode, err := usr.custRepository.ResendEmail(ctx, authData.Uid, p.Email)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	bodyEmail := "Hallo <b>" + authData.CustName + "</b>!, <br> Terimakasih telah bersedia bergabung bersama kami, silahkan lakukan verifikasi email anda dengan klik link berikut : " + CONFIG_API_URL + "/v1/verify?val=" + emailToken
	err = sendMail(p.Email, "Selamat Menjadi Bagian Pengguna Semesta Ban!", bodyEmail) // keep go

	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrSendEmail, crashy.Message(crashy.ErrSendEmail)), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)
}

func (usr *UsersHandler) RequestPinEmail(w http.ResponseWriter, r *http.Request) {
	var (
		p        ResendEmailRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	pin, errCode, err := usr.custRepository.RequestPinEmail(ctx, authData.Uid, p.Email)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	bodyEmail := "Hallo <b>" + authData.CustName + "</b>!<br> Anda telah melakukan request untuk pergantian email, berikut adalah kode yang dibutuhkan unik untuk diinput kedalam aplikasi untuk mengganti email anda : " + pin
	err = sendMail(p.Email, "Selamat Menjadi Bagian Pengguna Semesta Ban!", bodyEmail)

	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrSendEmail, crashy.Message(crashy.ErrSendEmail)), http.StatusInternalServerError)
		return
	}

	// fmt.Println(pin)

	response.Yay(w, r, "success", http.StatusOK)
}

func (usr *UsersHandler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	var (
		p                ChangeEmailRequest
		ctx              = r.Context()
		authData         = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		hashedTokenEmail = helper.GenerateHashString()
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	errCode, err := usr.custRepository.ChangeEmail(ctx, authData.Uid, p.OldEmail, p.NewEmail, hashedTokenEmail, p.Code)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//temporary
	bodyEmail := "Hallo <b>" + authData.CustName + "</b>!, <br> Terimakasih telah bersedia bergabung bersama kami, silahkan lakukan verifikasi email anda dengan klik link berikut : " + CONFIG_API_URL + "/v1/verify?val=" + hashedTokenEmail
	_ = sendMail(p.NewEmail, "Selamat Menjadi Bagian Pengguna Semesta Ban!", bodyEmail) // keep going even though send email failed

	response.Yay(w, r, "success", http.StatusOK)
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
