package repo_customers

import "context"

type CustomersRepository interface {
	Login(ctx context.Context, email, password string) (res Customers, errCode string, err error)
	CheckEmailExist(ctx context.Context, email string) (res bool, errCode string, err error)
	Register(ctx context.Context, name, email, emailToken, password, uid string) (errCode string, err error)
	VerifyEmail(ctx context.Context, emailToken string) (errCode string, err error)
	GetCustomer(ctx context.Context, uid string) (res Customers, errCode string, err error)
	ChangePassword(ctx context.Context, uid, oldPass, newPass string) (errCode string, err error)
	ResendEmail(ctx context.Context, uid, email string) (emailToken, errCode string, err error)
	RequestPinEmail(ctx context.Context, uid, email string) (pin, errCode string, err error)
	ChangeEmail(ctx context.Context, uid, oldEmail, newEmail, hashedTokenEmail, code string) (errCode string, err error)
}
