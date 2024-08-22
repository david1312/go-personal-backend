package constants

import "time"

const (
	OrderAsc  = "Asc"
	OrderDesc = "Desc"

	//
	StatusSuccess = "success"
	StatusFailed  = "failed"

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

	TransactionSourceApp = "APP"

	CSNumber = "081217950269"

	TransStatusBerhasil = "Berhasil"

	StaticPromoBanner = "banner-promo-agustus-2022.png"

	DefaultImgPng = "default.png"

	TF_BNI     = "TF_BNI"
	TF_PERMATA = "TF_PERMATA"
	TF_BRI     = "TF_BRI"
	COD        = "COD"

	TranStatusPembayaran = "Menunggu Pembayaran"
	TranStatusDipasang   = "Menunggu Dipasang"
	TranStatusDiproses   = "Diproses"
	TranStatusBerhasil   = "Berhasil"
	TranStatusBatal      = "Pesanan Dibatalkan"
	TranStatusSelesai    = "Selesai"

	DirectPayment = "Bayar Langsung"

	DIR_FILES        = "files/"
	DIR_REPORT_SALES = "reports/"
	MAX_COMMON_SIZE  = 10

	FORMAT_EXCEL = "sales-*.xlsx"

	// register related
	MinimumLengthPassword = 6
	// token related
	LoginTokenExpiry     = 24 * 7 * time.Hour  // 7 days
	RefreshTokenExpiry   = 24 * 14 * time.Hour // 14 days
	AnonymousTokenExpiry = 24 * 30 * time.Hour // 30 days
	ApiVersion           = "v1.0.0"

	// others
	BannerMsgToday    = "Ban pilihan mu akan dipasang hari ini"
	BannerMsgUpcoming = "Ban pilihan mu akan dipasang dalam %v hari"
)

var (
	ScheduleTime               []string = []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00", "19:00"}
	TransStatus                []string = []string{TranStatusPembayaran, TranStatusDipasang, TranStatusDiproses, TranStatusBerhasil, TranStatusBatal, TranStatusSelesai}
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

	FEE_LAZADA    = 8.75
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

var (
	JU_USER    = "davidbernadi13@gmail.com"
	JU_PASS    = "Mitoma13@@"
	JU_EXPIRED = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IlVTRVI6ZGF2aWRiZXJuYWRpMTNAZ21haWwuY29tOjE3Mi43MC4xODguNDYiLCJleHAiOjEwNjk2NDk1MzQsImlzX3dtc19taWdyYXRlZCI6dHJ1ZSwiaWF0IjoxNjgxMDI2NDQ5fQ.aGK9Xz4fD3fsRnhG6yopb1NhPLluJlakAN2pTsQ9Xxk"
)

// push notif fcm related
var (
	PushNotifAction     = "page_transaction"
	PushNotifMsgSuccess = "Selamat pembayaran untuk invoice %v berhasil diterima, silahkan datang ke outlet kami dijadwal yang sudah anda pilih untuk mendapatkan free pemasangan ban"
	PushNotifTitle      = "Pembayaran Berhasil Diterima"
)
