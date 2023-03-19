package repo_customers

import "context"

type CustomersRepository interface {
	Login(ctx context.Context, email, password string) (res Customers, errCode string, err error)
	UpdateDeviceToken(ctx context.Context, email, deviceToken string) (errCode string, err error) //todo handle multiple device?
	CheckEmailExist(ctx context.Context, email string) (res bool, errCode string, err error)
	CheckPhoneExist(ctx context.Context, phone string) (res bool, errCode string, err error)
	Register(ctx context.Context, name, email, emailToken, password, uid, phone string) (cleanUid, errCode string, err error)
	RegisterFromGoogleSignin(ctx context.Context, name, email, uid, avatar string) (cleanUid, errCode string, err error)
	VerifyEmail(ctx context.Context, emailToken string) (errCode string, err error)
	GetCustomer(ctx context.Context, uid string) (res Customers, errCode string, err error)
	GetCustomerByEmail(ctx context.Context, email string) (res Customers, errCode string, err error)
	ChangePassword(ctx context.Context, uid, oldPass, newPass string) (errCode string, err error)
	ResendEmail(ctx context.Context, uid, email string) (emailToken, errCode string, err error)
	RequestPinEmail(ctx context.Context, uid, email string) (pin, errCode string, err error)
	ChangeEmail(ctx context.Context, uid, oldEmail, newEmail, hashedTokenEmail, code string) (errCode string, err error)
	UpdateName(ctx context.Context, uid, name string) (errCode string, err error)
	UpdatePhoneNumber(ctx context.Context, uid, phone string) (errCode string, err error)
	UpdateGender(ctx context.Context, uid, gender string) (errCode string, err error)
	UpdateBirthDate(ctx context.Context, uid, birthdate string) (errCode string, err error)
	UploadProfileImg(ctx context.Context, uid, imgName string) (errCode string, err error)
}
