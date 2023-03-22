package repo_master_data

import "context"

type MasterDataRepository interface {
	Magic(ctx context.Context) error
	UpdateTransactionExpired(ctx context.Context) error

	GetListMerkBan(ctx context.Context) (res []MerkBan, errCode string, err error)
	GetListRingBan(ctx context.Context) (res []string, errCode string, err error)
	GetListUkuranBan(ctx context.Context) (res []UkuranRingBan, errCode string, err error)
	GetListUkuranBanByBrandMotor(ctx context.Context, idBrandMotor []int) (res []UkuranRingBan, errCode string, err error)
	GetListUkuranBanByMotor(ctx context.Context, idMotor int) (res []UkuranRingBan, errCode string, err error)
	GetListMerkMotor(ctx context.Context) (res []MerkMotor, errCode string, err error)
	GetListMotorByBrand(ctx context.Context, idBrandMotor int) (res []Motor, errCode string, err error)
	GetListPaymentMethod(ctx context.Context) (res []PaymentMethod, errCode string, err error)
	GetListTopRankpMotor(ctx context.Context) (res []Motor, errCode string, err error)
	GetListUkuranBanRaw(ctx context.Context) (res []UkuranRingBan, errCode string, err error)

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

	//motor related
	GetListMotor(ctx context.Context, fp ListMotorRequestRepo) (res []MotorMD, totalData int, errCode string, err error)
	MotorAdd(ctx context.Context, name, idBrandMotor, idCategoryMotor, icon string) (errCode string, err error)
	MotorUpdate(ctx context.Context, idMotor int, name string, idBrandMotor, idCategoryMotor int) (errCode string, err error)
	MotorCheckExists(ctx context.Context, idMotor string) (exists bool, errCode string, err error)
	MotorCheckUsed(ctx context.Context, idMotor string) (exists bool, errCode string, err error)
	MotorUpdateImage(ctx context.Context, idMotor, fileName, uploadPath, dirFile string) (errCode string, err error)
	MotorRemove(ctx context.Context, idMotor, uploadPath, dirFile string) (errCode string, err error)

	//tire size
	TireSizeExist(ctx context.Context, id string) (exists bool, errCode string, err error)
	TireSizeAdd(ctx context.Context, id, idRing, idSize string) (errCode string, err error)
	TireSizeUsed(ctx context.Context, id string) (exists bool, errCode string, err error)
	TireSizeDelete(ctx context.Context, id string) (errCode string, err error)

	//tire ring
	TireRingExist(ctx context.Context, id int) (exists bool, errCode string, err error)
	TireRingAdd(ctx context.Context, id int, nameRing string) (errCode string, err error)

	GetListCategoryMotor(ctx context.Context) (res []CategoryMotor, errCode string, err error)
}
