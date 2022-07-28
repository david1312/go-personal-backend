package helper

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"semesta-ban/pkg/constants"
	"strconv"
	"strings"
	"time"

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
