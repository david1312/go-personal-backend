package constants

const (
	OrderAsc  = "ASC"
	OrderDesc = "Desc"

	UserDir           = "usr/"
	ProductDir        = "product/"
	TireBrandDir      = "merkban/"
	MotorBrandDir     = "merkmotor/"
	MotorDir          = "motor/"
	PaymentMethod     = "payment_method/"
	RatingsProductDir = "rating_product/"
	RatingOutletDir   = "rating_outlet/"
	MerchantDir       = "merchant/"
	PromoDir          = "promo/"
	MotorCategoryDir  = "motor_category/"

	Male        = "Laki-laki"
	Female      = "Perempuan"
	OtherGender = "Yang Lain"

	TrTypeBooking = "Booking Outlet"
	TrTypeKirim   = "Kirim Barang"

	LimitOrderInHour = 3

	FraudAccept    = "accept"
	FraudDeny      = "deny"
	FraudChallenge = "challenge"

	PaymentBankTransfer = "bank_transfer"
	BankBNI             = "bni"
	BankPermata         = "bni"
	BankBri             = "bri"

	Bank
	DBPaymentSettle    = "Lunas"
	DBPaymentNotSettle = "Belum Lunas"

	MTransStatusSettlement = "settlement"
	MTransStatusPending    = "pending"
	MTransStatusExpire     = "expire"
	MTransStatusCancel     = "cancel"
	MTransStatusDeny       = "deny"

	CSNumber = "081217950269"

	TransStatusBerhasil = "Berhasil"

	StaticPromoBanner = "banner-promo-agustus-2022.png"

	DefaultImgPng = "default.png"

	TF_BNI     = "TF_BNI"
	TF_PERMATA = "TF_PERMATA"
	TF_BRI     = "TF_BRI"
	COD        = "COD"

	T_STATUS_MENUNGGU_DIPASANG = "Menunggu Dipasang"
)

var (
	ScheduleTime []string = []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00", "19:00"}
	TransStatus  []string = []string{"Menunggu Pembayaran", "Menunggu Dipasang", "Diproses", "Berhasil", "Pesanan Dibatalkan", "Selesai"}
)
