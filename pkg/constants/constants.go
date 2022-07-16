package constants

const (
	OrderAsc  = "ASC"
	OrderDesc = "Desc"

	UserDir       = "usr/"
	ProductDir    = "product/"
	TireBrandDir  = "merkban/"
	MotorBrandDir = "merkmotor/"
	MotorDir      = "motor/"
	PaymentMethod = "payment_method/"

	Male        = "Laki-laki"
	Female      = "Perempuan"
	OtherGender = "Yang Lain"

	TrTypeBooking = "Booking Outlet"
	TrTypeKirim   = "Kirim Barang"

	LimitOrderInHour = 3
)

var (
	ScheduleTime []string = []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00", "19:00"}
	TransStatus  []string = []string{"Menunggu Pembayaran", "Menunggu Konfirmasi", "Menunggu Kedatangan", "Diproses", "Selesai"}
)
