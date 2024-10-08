package crashy

const (
	//ErrCodeUnexpected code for generic error for unrecognized cause
	ErrCodeUnexpected = "ERR_UNEXPECTED"
	//ErrCodeNetBuild code for resource connection build issue
	ErrCodeNetBuild = "ERR_NET_BUILD"
	//ErrCodeNetConnect code for resource connection establish issue
	ErrCodeNetConnect = "ERR_NET_CONNECT"
	//ErrCodeValidation code for any data validation issues
	ErrCodeValidation = "ERR_VALIDATION"
	//ErrCodeFormatting code for any formatting issue(s) includes marshalling and unmarshalling
	ErrCodeFormatting = "ERR_FORMATTING"
	//ErrCodeDataRead code for any storage read issue(s)
	ErrCodeDataRead = "ERR_DATA_READ"
	//ErrCodeDataWrite code for any storage write issue(s)
	ErrCodeDataWrite = "ERR_DATA_WRITE"
	//ErrCodeNoResult code when data provider given no result for any query
	ErrCodeNoResult = "ERR_NO_RESULT"
	//ErrCodeUnauthorized code when any access doesnt contains enough authorization
	ErrCodeUnauthorized = "ERR_UNAUTHORIZED"
	ErrCodeExpired      = "ERR_EXPIRED"
	ErrCodeForbidden    = "ERR_FORBIDDEN"

	//ErrCodeTooManyRequest code when a user token used more than given rate
	ErrCodeTooManyRequest = "ERR_REQUEST_LIMIT"
	ErrCodeDataIncomplete = "ERR_DATA_INCOMPLETE"

	//ErrCodeEncryptData code for any encrypting data issue(s)
	ErrCodeEncryptData = "ERR_ENCRYPT_DATA"

	//Err code when failed get activity data
	ErrCodeGetActivtiy = "ERR_DATA_GET_ACTIVITY"
	//Err code when failed get template data
	ErrCodeGetTemplate = "ERR_DATA_GET_TEMPLATE"
	//Err code when resend time abusing the specified time / config
	ErrCodeTimeResend = "ERR_DATA_TIME_RESEND"

	ErrInvalidToken = "ERR_INVALID_TOKEN"
	ErrServer
	ErrExceededFileSize = "ERR_EXCEEDED_FILE_SIZE"
	ErrFileNotFound     = "ERR_FILE_NOT_FOUND"
	ErrUploadFile       = "ERR_UPLOAD_FILE"

	//Err code Login related
	ErrInvalidUser = "ERR_INVALID_USER"

	//Err code Register related
	ErrEmailExists   = "ERR_EMAIL_EXISTS"
	ErrShortPassword = "ERR_SHORT_PWD"
	ErrPhoneExists   = "ERR_PHONE_EXISTS"

	ErrSendEmail         = "ERR_SEND_EMAIL"
	ErrInvalidTokenEmail = "ERR_INVALID_TOKEN_EMAIL"

	ErrInvalidOldPassword = "ERR_INVALID_OLD_PWD"
	ErrSamePassword       = "ERR_SAME_PWD"
	ErrInvalidGender      = "ERR_INVALID_GENDER"
	ErrInvalidEmail       = "ERR_INVALID_EMAIL"
	ErrInvalidCode        = "ERR_INVALID_CODE"
	ErrInvalidPhone       = "ERR_INVALID_PHONE"
	ErrSamePhoneSelf      = "ERR_SAME_PHONE_SELF"
	ErrDedupPhone         = "ERR_DEDUP_PHONE"
	ErrInvalidBirthDate   = "ERR_INVALID_BIRTHDATE"

	//Err code product related
	ErrInvalidProductID = "ERR_INVALID_PRODUCT_ID"
	ErrExistsWishlist   = "ERR_EXIST_WISHLIST"

	//err code transactions
	ErrBackwardDate        = "ERR_BACKWARD_DATE"
	ErrInvalidScheduleTime = "ERR_INVALID_SCHEDULE_TIME"
	ErrInsufficientStock   = "ERR_INSUFFICIENT_STOCK"

	//midtrans error
	ErrRequestMidtrans = "ERR_REQUEST_MIDTRANS"

	//merchant
	ErrInvalidUserMerchant = "ERR_INVALID_USER_MERCHANT"
	ErrBrandMotorUsed      = "ERR_BRAND_MOTOR_USED"
	ErrInvalidBrandMotor   = "ERR_INVALID_BRAND_MOTOR"
	ErrTireBrandUsed       = "ERR_TIRE_BRAND_USED"
	ErrInvalidTireBrand    = "ERR_INVALID_TIRE_BRAND"
	ErrCommonInvalid       = "ERR_INVALID_COMMON"
	ErrMotorUsed           = "ERR_MOTOR_USED"

	ErrMinimumPhoto = "ERR_MINIMUM_PHOTO"

	ErrTireSizeExists = "ERR_TIRE_SIZE_EXISTS"
	ErrTireSizeUsed   = "ERR_TIRE_SIZE_USED"

	ErrTireRingExists = "ERR_TIRE_RING_EXISTS"

	ErrInvalidLocation = "ERR_INVALID_LOCATION"
)

var mapper = map[ErrCode]string{
	ErrInvalidToken:        "error veryfying token",
	ErrCodeUnexpected:      "maaf, terjadi gangguan pada server",
	ErrCodeNetBuild:        "failed to build connection to data source",
	ErrCodeNetConnect:      "failed to establish connection to data source",
	ErrCodeValidation:      "request contains invalid data",
	ErrCodeFormatting:      "an error occurred while formatting data",
	ErrCodeDataRead:        "maaf, terjadi gangguan pada server",
	ErrCodeDataWrite:       "failed to persist data into provider",
	ErrCodeNoResult:        "no result found match criteria",
	ErrCodeUnauthorized:    "unauthorized access",
	ErrCodeForbidden:       "forbidden access",
	ErrCodeExpired:         "expired pemission",
	ErrCodeTooManyRequest:  "request limit exceeded",
	ErrCodeDataIncomplete:  "stored data incomplete",
	ErrCodeEncryptData:     "failed to encrypting data",
	ErrInvalidUser:         "email atau password yang anda masukan salah",
	ErrEmailExists:         "maaf email yang anda masukan sudah terdaftar",
	ErrPhoneExists:         "maaf nomor handphone yang anda masukan sudah terdaftar",
	ErrSendEmail:           "terjadi kesalahan saat mengirim email, mohon coba beberapa saat lagi",
	ErrInvalidTokenEmail:   "token untuk verifikasi email tidak valid",
	ErrInvalidOldPassword:  "password lama yang anda masukan salah",
	ErrShortPassword:       "password minimum 6 karakter",
	ErrSamePassword:        "password baru tidak boleh sama dengan password lama",
	ErrInvalidEmail:        "email yang anda masukan tidak sesuai",
	ErrInvalidCode:         "kode yang anda masukan salah",
	ErrInvalidProductID:    "product id yang anda masukan salah",
	ErrInvalidGender:       "value gender yang anda masukan tidak sesuai format",
	ErrInvalidPhone:        "nomor hp yang anda masukan tidak sesuai format",
	ErrSamePhoneSelf:       "tidak bisa mengganti dengan nomor hp yang sama dengan nomor sebelumnya",
	ErrDedupPhone:          "nomor hp yang anda masukan sudah digunakan",
	ErrInvalidBirthDate:    "format tanggal wajib yyyy-mm-dd (contoh : 1996-12-02)",
	ErrExistsWishlist:      "produk yang dipilih sudah berada dalam wishlist anda",
	ErrExceededFileSize:    "maximum file/foto yang di upload adalah : ",
	ErrFileNotFound:        "file tidak ditemukan",
	ErrUploadFile:          "terjadi gangguan saat upload file ke server, mohon coba beberapa saat lagi",
	ErrBackwardDate:        "tidak bisa melakukan pemesanan di hari yang sudah lewat",
	ErrInvalidScheduleTime: "format schedule_time tidak sesuai",
	ErrInsufficientStock:   "transaksi gagal, stock barang tidak mencukupi",
	ErrRequestMidtrans:     "transaksi berhasil terbuat namun terdapat gangguan saat melakukan request payment ke midtrans, segera hubungi admin",
	ErrBrandMotorUsed:      "Mohon maaf, brand yang ingin dihapus sedang dipakai oleh motor tertentu. Mohon pastikan brand yang ingin dihapus tidak terikat dengan motor manapun",
	ErrInvalidBrandMotor:   "Id merk motor yang anda input tidak terdaftar didalam sistem",
	ErrTireBrandUsed:       "Mohon maaf, merk ban yang ingin dihapus masih terpakai didalam produk tertentu. Mohon pastikan merk ban yang ingin dihapus tidak terikat dengan produk manapun",
	ErrInvalidTireBrand:    "Id merk ban yang anda input tidak terdaftar didalam sistem",
	ErrCommonInvalid:       "Data yang diinput tidak ditemukan",
	ErrMotorUsed:           "Id Motor yang diinput telah terpakai dalam matriks relasi motor dengan ukuran ban. Mohon pastikan id motor yang ingin dihapus tidak sedang dipakai oleh matriks tersebut.",

	//merchant
	ErrInvalidUserMerchant: "username atau password yang anda masukan salah",
	ErrMinimumPhoto:        "1 Produk minimal memiliki 1 gambar / foto barang",

	ErrTireSizeExists: "Ukuran ban yang ingin ditambahkan sudah terdaftar",
	ErrTireSizeUsed:   "Hapus gagal, ukuran ban yang ingin dihapus sedang dipakai dalam matriks relasi motor dengan ukuran ban atau dipakai oleh produk tertentu.",
	ErrTireRingExists: "Ukuran ring ban yang ingin ditambahkan sudah terdaftar",

	ErrInvalidLocation: "Maaf anda berada diluar jangkauan outlet kami, nantikan pembukaan outlet kami di kota anda segera",
}

// Message retrieve error messages from given error code
func Message(code ErrCode) string {
	if s, ok := mapper[code]; ok {
		return s
	}
	return mapper[ErrCodeUnexpected]
}

// Messages retrieve all registered mapping error messages
func Messages() map[ErrCode]string {
	return mapper
}
