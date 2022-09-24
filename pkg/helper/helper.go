package helper

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/log"
	"strconv"
	"strings"
	"time"

	b64 "encoding/base64"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func RandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func RandomNumber(n int) string {
	const alphanum = "0123456789"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func GenerateHashString() string {
	h := sha256.New()
	h.Write([]byte(RandomString(6)))
	hashedStr := hex.EncodeToString(h.Sum(nil))

	return hashedStr
}

func IsFloatNoDecimal(val float64) bool {
	return val == float64(int(val))
}

func ValidateParam(val string) bool {
	valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
	return valid.MatchString(val)
}

func IsStringNumeric(val string) bool {
	valid := regexp.MustCompile("^[0-9_]+$")
	return valid.MatchString(val)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ConvertPhoneNumber(inputNumber string) (res string) {
	res = inputNumber
	if res[:1] == "+" {
		res = strings.Replace(res, "+", "", 1)
	}
	if res[:1] == "0" {
		res = strings.Replace(res, "0", "62", 1)
	}
	res = strings.ReplaceAll(res, " ", "")
	return
}

func GenerateCustomerId(lastId string) (res string) {
	splitted := strings.Split(lastId, "-")

	test, _ := strconv.Atoi(splitted[1])
	test++
	secRes := strconv.Itoa(test)
	if test <= 9 {
		secRes = fmt.Sprintf("000%v", test)
	} else if test <= 99 && test > 9 {
		secRes = fmt.Sprintf("00%v", test)
	} else if test <= 999 && test > 99 {
		secRes = fmt.Sprintf("0%v", test)
	}

	res = fmt.Sprintf("%s-%v", splitted[0], secRes)
	return
}

func GenerateTransactionId(lastId, date string) (res string) {
	if len(lastId) == 0 {
		res = fmt.Sprintf("INV-%s-0001", date)
		return
	}
	splitted := strings.Split(lastId, "-")

	test, _ := strconv.Atoi(splitted[2])
	test++
	secRes := strconv.Itoa(test)
	if test <= 9 {
		secRes = fmt.Sprintf("000%v", test)
	} else if test <= 99 && test > 9 {
		secRes = fmt.Sprintf("00%v", test)
	} else if test <= 999 && test > 99 {
		secRes = fmt.Sprintf("0%v", test)
	}

	res = fmt.Sprintf("%s-%s-%v", splitted[0], date, secRes)
	return
}

func ConvertFileSizeToMb(size int) (res int) {
	return size * 1000000
}

func GetUploadedFileName(file string) string {
	splitted := strings.Split(file, "/")
	return splitted[4]
}

func ValidateScheduleTime(schedule string) bool {
	set := make(map[string]bool)
	for _, v := range constants.ScheduleTime {
		set[v] = true
	}

	return (set[schedule])
}

func StringInArray(arr []string, param string) bool {
	set := make(map[string]bool)
	for _, v := range arr {
		set[v] = true
	}

	return (set[param])
}

func FormatCurrency(number int) string {
	p := message.NewPrinter(language.English)
	test := p.Sprintf("%d", number)
	return "Rp" + strings.ReplaceAll(test, ",", ".")
}

func CreateHttpClient(ctx context.Context, timeout int, skipSSL bool) *http.Client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
	}

	return &http.Client{
		Timeout:   time.Duration(timeout) * time.Minute,
		Transport: tr,
	}
}

func GenerateB64AuthMidtrans(serverKey string) string {
	if len(serverKey) == 0 {
		return ""
	}
	return b64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
}

func FormatInstallationTime(dateStr, timeStr string) string {
	date, _ := time.Parse("2006-01-02", dateStr[:10])

	dayName := MappingDaysNameId(date.Weekday().String())
	return fmt.Sprintf("%v, %v, %v WIB", dayName, date.Format("02 January 2006"), timeStr)
}

func MappingBankName(paymentMethod string) string {
	switch paymentMethod {
	case "TF_BNI":
		return "BNI"
	case "TF_BCA":
		return "BCA"
	case "TF_BRI":
		return "BRI"
	case "TF_MANDIRI":
		return "MANDIRI"
	case "TF_PERMATA":
		return "PERMATA"
	case "TF_GOPAY":
		return "GOPAY"
	default:
		return ""
	}

}

func MappingDaysNameId(dayName string) string {
	switch dayName {
	case "Sunday":
		return "Minggu"
	case "Monday":
		return "Senin"
	case "Tuesday":
		return "Selasa"
	case "Wednesday":
		return "Rabu"
	case "Thursday":
		return "Kamis"
	case "Friday":
		return "Jumat"
	case "Saturday":
		return "Sabtu"
	default:
		return ""
	}
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func UploadSingleImage(r *http.Request, fieldName, uploadPath, directory string, maxSize int) (fileName, errCode string, err error) {
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(fieldName)
	if err != nil {
		errCode = crashy.ErrFileNotFound
		return
	}
	defer file.Close()

	if handler.Size > int64(ConvertFileSizeToMb(maxSize)) {
		errCode = crashy.ErrExceededFileSize
		return
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(uploadPath+directory, "pic-*.png")
	if err != nil {
		errCode = crashy.ErrUploadFile
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		errCode = crashy.ErrUploadFile
		return
	}
	// write this byte array to our temporary file
	fileName = GetUploadedFileName(tempFile.Name())

	tempFile.Write(fileBytes)
	tempFile.Chmod(0604)
	log.Infof("success upload %s to the server x \n", fileName)
	return
}

func UploadImage(r *http.Request, fieldName, uploadPath, directory string) (fileNameList []string, errCode string, err error) {
	for _, fh := range r.MultipartForm.File[fieldName] {
		f, errTemp := fh.Open()
		if err != nil {
			// Handle error
			err = errTemp
			errCode = crashy.ErrFileNotFound
			break
		}

		tempFile, errTemp := ioutil.TempFile(uploadPath+directory, "pic-*.png")
		if err != nil {
			err = errTemp
			errCode = crashy.ErrUploadFile
			break
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, errTemp := ioutil.ReadAll(f)
		if err != nil {
			err = errTemp
			errCode = crashy.ErrUploadFile
			break
		}
		// write this byte array to our temporary file
		fileName := GetUploadedFileName(tempFile.Name())

		tempFile.Write(fileBytes)
		tempFile.Chmod(0604)
		log.Infof("success upload %s to the server x \n", fileName)
		fileNameList = append(fileNameList, fileName)

		// Read data from f
		f.Close()
	}
	return

}
