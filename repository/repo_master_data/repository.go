package repo_master_data

import "context"

type MasterDataRepository interface {
	GetListMerkBan(ctx context.Context) (res []MerkBan, errCode string, err error)
	GetListUkuranBan(ctx context.Context) (res []UkuranRingBan, errCode string, err error)
	GetListUkuranBanByBrandMotor(ctx context.Context, idBrandMotor []int) (res []UkuranRingBan, errCode string, err error)
	GetListUkuranBanByMotor(ctx context.Context, idMotor int) (res []UkuranRingBan, errCode string, err error)
	GetListMerkMotor(ctx context.Context) (res []MerkMotor, errCode string, err error)
	GetListMotorByBrand(ctx context.Context, idBrandMotor int) (res []Motor, errCode string, err error)
	GetListPaymentMethod(ctx context.Context) (res []PaymentMethod, errCode string, err error)
}
