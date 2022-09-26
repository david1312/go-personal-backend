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
	GetListTopRankpMotor(ctx context.Context) (res []Motor, errCode string, err error)

	//merk motor related
	AddBrandMotor(ctx context.Context, name, icon string) (errCode string, err error)
	CheckBrandMotorUsed(ctx context.Context, idMotor int) (exists bool, errCode string, err error)
	RemoveBrandMotor(ctx context.Context, idMotor int, uploadPath, dirFile string) (errCode string, err error)
	UpdateBrandMotor(ctx context.Context, idMotor int, name string) (errCode string, err error)
	UpdateBrandMotorImage(ctx context.Context, idMotor int, fileName, uploadPath, dirFile string) (errCode string, err error)
	CheckBrandMotorExist(ctx context.Context, idMotor int) (exists bool, errCode string, err error)

	//merk ban related
	AddTireBrand(ctx context.Context, id, name, icon, ranking string) (errCode string, err error)
	CheckTireBrandUsed(ctx context.Context, idMerkBan string) (exists bool, errCode string, err error)
	RemoveTireBrand(ctx context.Context, idMerkBan, uploadPath, dirFile string) (errCode string, err error)
	UpdateTireBrand(ctx context.Context, idMerkBan, name string, ranking int) (errCode string, err error)
	CheckTireBrandExist(ctx context.Context, idMerkBan string) (exists bool, errCode string, err error)
	UpdateTireBrandImage(ctx context.Context, idMerkBan, fileName, uploadPath, dirFile string) (errCode string, err error)
}
