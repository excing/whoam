package main

import (
	"math/rand"
	"regexp"
	"time"
	"unsafe"
)

// VerifyEmailFormat Verify email address
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` // email regular expression
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

const digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_."

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// see: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go?answertab=votes#tab-top
func genRandCode(n int, bit ...int) string {
	b := make([]byte, n)
	start, end := 0, len(digits)
	switch len(bit) {
	case 1:
		end = bit[0]
	case 2:
		start = bit[0]
		end = bit[1]
	}
	dict := digits[start:end]

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(dict) {
			b[i] = dict[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// New64BitUUID 获取 64进制 UUID, 10个数字+26个小写字母+26个大写字母+(=、_), 长度 24
// generates a random UUID according to RFC 4122
func New64BitUUID() (string, error) {
	return genRandCode(64), nil
}
