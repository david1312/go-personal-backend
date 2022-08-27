package repo_merchant

import "context"

type MerchantRepository interface {
	Login(ctx context.Context, username, password string) (res MerchantData, errCode string, err error)
	GetMerchantProfile(ctx context.Context, username string) (res MerchantData, errCode string, err error)
}
