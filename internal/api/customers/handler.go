package customers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/pkg/log"
	custRepo "semesta-ban/repository/repo_customers"
	"strings"
	"time"

	localMdl "semesta-ban/internal/api/middleware"
	cn "semesta-ban/pkg/constants"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"gopkg.in/gomail.v2"
)

// todo move to config
const CONFIG_SMTP_HOST = "mail.sunmorisemestaban.com"
const CONFIG_SMTP_PORT = 465
const CONFIG_SENDER_NAME = "PT. Sunmori Semesta Ban <support@sunmorisemestaban.com>"
const CONFIG_AUTH_EMAIL = "support@sunmorisemestaban.com"
const CONFIG_AUTH_PASSWORD = "spyxfamily13"
const CONFIG_API_URL = "https://api.sunmorisemestaban.com"

type UsersHandler struct {
	db                *sqlx.DB
	custRepository    custRepo.CustomersRepository
	jwt               *localMdl.JWT
	baseAssetUrl      string
	uploadPath        string
	profilePicPath    string
	profilePicMaxSize int
}

//todo REMEMBER 30 May gmail tidak support lagi less secure app find solution

func NewUsersHandler(db *sqlx.DB, cr custRepo.CustomersRepository, jwt *localMdl.JWT, baseAssetUrl, uploadPath,
	profilePicPath string, profilePicMaxSize int) *UsersHandler {
	return &UsersHandler{db: db, custRepository: cr, jwt: jwt, baseAssetUrl: baseAssetUrl, uploadPath: uploadPath, profilePicPath: profilePicPath, profilePicMaxSize: profilePicMaxSize}
}

func (usr *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		p        RegisterRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
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

	cleanUid, errCode, err := usr.custRepository.Register(ctx, p.Name, p.Email, hashedTokenEmail, p.Password, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	//temporary
	bodyEmail := "Hallo <b>" + p.Name + "</b>!, <br> Terimakasih telah bersedia bergabung bersama kami, silahkan lakukan verifikasi email anda dengan klik link berikut : " + CONFIG_API_URL + "/v1/verify?val=" + hashedTokenEmail
	_ = sendMail(p.Email, "Selamat Menjadi Bagian Pengguna Semesta Ban!", bodyEmail) // keep going even though send email failed

	//generate token
	expiredTime := time.Now().Add(3 * time.Minute)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      cleanUid,
		CustName: p.Name,

		Expired: expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Minute * 4)
	_, tokenRefresh, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      cleanUid,
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
	expiredTime := time.Now().Add(3 * time.Minute)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      customer.Uid,
		CustName: customer.Name,
		Expired:  expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Minute * 4 )
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
	avatar := ""

	if len(customer.Avatar.String) > 0 && customer.Avatar.String[:3] == "pic" {
		avatar = usr.baseAssetUrl + usr.profilePicPath + customer.Avatar.String
	} else if len(customer.Avatar.String) > 0 && customer.Avatar.String[:3] != "pic" {
		avatar = customer.Avatar.String
	}

	response.Yay(w, r, GetCustomerResponse{
		Name:          customer.Name,
		Email:         customer.Email,
		EmailVerified: isEmailVerified,
		Phone:         customer.Phone.String,
		PhoneVerified: isPhoneVerified,
		Gender:        customer.Gender.String,
		Avatar:        avatar,
		Birthdate:     birthdateVal,
		CustId:        customer.CustId,
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

func (usr *UsersHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	var (
		p        UpdateNameRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}
	errCode, err := usr.custRepository.UpdateName(ctx, authData.Uid, p.Name)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (usr *UsersHandler) UpdatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	var (
		p        UpdatePhoneRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	if !helper.IsStringNumeric(p.Phone) {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidPhone), crashy.ErrInvalidPhone, crashy.Message(crashy.ErrCode(crashy.ErrInvalidPhone))), http.StatusBadRequest)
		return
	}

	errCode, err := usr.custRepository.UpdatePhoneNumber(ctx, authData.Uid, helper.ConvertPhoneNumber(p.Phone))
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (usr *UsersHandler) UpdateGender(w http.ResponseWriter, r *http.Request) {
	var (
		p        UpdateGenderRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	checkValidGender := helper.StringInSlice(p.Gender, []string{cn.Male, cn.Female, cn.OtherGender})

	if !checkValidGender {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidGender), crashy.ErrInvalidGender, crashy.Message(crashy.ErrCode(crashy.ErrInvalidGender))), http.StatusBadRequest)
		return
	}
	errCode, err := usr.custRepository.UpdateGender(ctx, authData.Uid, p.Gender)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (usr *UsersHandler) UpdateBirthDate(w http.ResponseWriter, r *http.Request) {
	var (
		p        UpdateBirthDateRequest
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	date := strings.Split(p.Birthdate, "-")
	if len(date) != 3 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidBirthDate), crashy.ErrInvalidBirthDate, crashy.Message(crashy.ErrCode(crashy.ErrInvalidBirthDate))), http.StatusBadRequest)
		return
	}

	if len(date[0]) != 4 || len(date[1]) != 2 || len(date[2]) != 2 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrInvalidBirthDate), crashy.ErrInvalidBirthDate, crashy.Message(crashy.ErrCode(crashy.ErrInvalidBirthDate))), http.StatusBadRequest)
		return
	}

	errCode, err := usr.custRepository.UpdateBirthDate(ctx, authData.Uid, p.Birthdate)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (usr *UsersHandler) UploadProfileImg(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
	)

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("profile_img")
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrFileNotFound, crashy.Message(crashy.ErrCode(crashy.ErrFileNotFound))), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Size > int64(helper.ConvertFileSizeToMb(usr.profilePicMaxSize)) {
		errMsg := fmt.Sprintf("%s%v mb", crashy.Message(crashy.ErrCode(crashy.ErrExceededFileSize)), usr.profilePicMaxSize)
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrExceededFileSize), crashy.ErrExceededFileSize, errMsg), http.StatusBadRequest)
		return
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(usr.uploadPath+usr.profilePicPath, "pic-*.png")
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrUploadFile, crashy.Message(crashy.ErrCode(crashy.ErrUploadFile))), http.StatusBadRequest)
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrUploadFile, crashy.Message(crashy.ErrCode(crashy.ErrUploadFile))), http.StatusBadRequest)
		return
	}
	// write this byte array to our temporary file
	fileName := helper.GetUploadedFileName(tempFile.Name())

	tempFile.Write(fileBytes)
	tempFile.Chmod(0604)
	log.Infof("success upload %s to the server x \n", fileName)

	errCode, err := usr.custRepository.UploadProfileImg(ctx, authData.Uid, fileName)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (usr *UsersHandler) SignInGoogle(w http.ResponseWriter, r *http.Request) {
	var (
		p   SignInGoogleRequest
		ctx = r.Context()
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	customer, errCode, err := usr.custRepository.GetCustomerByEmail(ctx, p.Email)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	if len(customer.Uid) > 0 {
		//generate token
		expiredTime := time.Now().Add(24 * time.Hour)
		_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
			Uid:      customer.Uid,
			CustName: customer.Name,
			Expired:  expiredTime,
		})

		//generate refresh token
		expiredTimeRefresh := time.Now().Add(time.Hour * 24 * 30)
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
		return

	}
	authData := ctx.Value(localMdl.CtxKey).(localMdl.Token)

	cleanUid, errCode, err := usr.custRepository.RegisterFromGoogleSignin(ctx, p.DisplayName, p.Email, authData.Uid, p.PhotoUrl)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//generate token
	expiredTime := time.Now().Add(24 * time.Hour)
	_, tokenLogin, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      cleanUid,
		CustName: p.DisplayName,

		Expired: expiredTime,
	})

	//generate refresh token
	expiredTimeRefresh := time.Now().Add(time.Hour * 24 * 30)
	_, tokenRefresh, _ := usr.jwt.JWTAuth.Encode(&localMdl.Token{
		Uid:      cleanUid,
		CustName: p.DisplayName,
		Expired:  expiredTimeRefresh,
	})

	response.Yay(w, r, LoginResponse{
		Token:        tokenLogin,
		ExpiredAt:    expiredTime,
		RefreshToken: tokenRefresh,
		RTExpired:    expiredTimeRefresh,
	}, http.StatusOK)

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
