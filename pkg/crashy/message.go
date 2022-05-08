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

	//Err code Login related
	ErrInvalidUser = "ERR_INVALID_USER"

	//Err code Register related
	ErrEmailExists = "ERR_EMAIL_EXISTS"

	ErrSendEmail         = "ERR_SEND_EMAIL"
	ErrInvalidTokenEmail = "ERR_INVALID_TOKEN_EMAIL"
)

var mapper = map[ErrCode]string{
	ErrCodeUnexpected:     "maaf, terjadi gangguan pada server",
	ErrCodeNetBuild:       "failed to build connection to data source",
	ErrCodeNetConnect:     "failed to establish connection to data source",
	ErrCodeValidation:     "request contains invalid data",
	ErrCodeFormatting:     "an error occurred while formatting data",
	ErrCodeDataRead:       "failed to read data from data provider",
	ErrCodeDataWrite:      "failed to persist data into provider",
	ErrCodeNoResult:       "no result found match criteria",
	ErrCodeUnauthorized:   "unauthorized access",
	ErrCodeForbidden:      "forbidden access",
	ErrCodeExpired:        "expired pemission",
	ErrCodeTooManyRequest: "request limit exceeded",
	ErrCodeDataIncomplete: "stored data incomplete",
	ErrCodeEncryptData:    "failed to encrypting data",
	ErrInvalidUser:        "email atau password yang anda masukan salah",
	ErrEmailExists:        "maaf email yang anda masukan sudah terdaftar",
	ErrSendEmail:          "terjadi kesalahan saat mengirim email",
	ErrInvalidTokenEmail:  "verifikasi email gagal, token untuk verifikasi email tidak valid",
}

//Message retrieve error messages from given error code
func Message(code ErrCode) string {
	if s, ok := mapper[code]; ok {
		return s
	}
	return mapper[ErrCodeUnexpected]
}

//Messages retrieve all registered mapping error messages
func Messages() map[ErrCode]string {
	return mapper
}
