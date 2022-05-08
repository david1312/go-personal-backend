package repo_customers

import "context"

type CustomersRepository interface {
	Login(ctx context.Context, email, password string) (nama, errCode string, err error)
	CheckEmailExist(ctx context.Context, email string) (res bool, errCode string, err error)
	Register(ctx context.Context, name, email, emailToken, password string) (errCode string, err error)
	VerifyEmail(ctx context.Context, emailToken string) (errCode string, err error)
}
