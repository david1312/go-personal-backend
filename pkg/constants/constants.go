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

	DirectPayment = "Bayar Langsung"

	DIR_FILES        = "files/"
	DIR_REPORT_SALES = "reports/"
	MAX_COMMON_SIZE  = 10

	FORMAT_EXCEL = "sales-*.xlsx"
)

var (
	ScheduleTime               []string = []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00", "19:00"}
	TransStatus                []string = []string{"Menunggu Pembayaran", "Menunggu Dipasang", "Diproses", "Berhasil", "Pesanan Dibatalkan", "Selesai"}
	TypeDataNumberReportsSales []int    = []int{8, 9, 10, 11, 12, 14, 15, 16, 17, 18, 19}
)

// merchant fee in percentage
var (
	CHANNEL_LAZADA    = "LAZADA"
	CHANNEL_JUBELIO   = "JUBELIO-POS"
	CHANNEL_INTERNAL  = "INTERNAL"
	CHANNEL_SHOPEE    = "SHOPEE"
	CHANNEL_AKULAKU   = "AKULAKU"
	CHANNEL_TIKTOK    = "TIKTOK"
	CHANNEL_TOKOPEDIA = "TOKOPEDIA"

	FEE_LAZADA    = 8.0
	FEE_TOKOPEDIA = 3.1
	FEE_SHOPEE    = 8.7
	FEE_TIKTOK    = 1.0 // +IDR 2000 exception only tiktok
	FEE_AKULAKU   = 2.0

	STATUS_COMPLETED  = "COMPLETED"
	STATUS_RETURNED   = "RETURNED"
	STATUS_FAILED     = "FAILED"
	STATUS_SHIPPED    = "SHIPPED"
	STATUS_PROCESSING = "PROCESSING"
	STATUS_PAID       = "PAID"
)
