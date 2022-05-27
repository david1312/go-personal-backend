package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
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
