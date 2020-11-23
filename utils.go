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

const digits = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz_."

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandNdigMbitString returns a randomly generated string with n digits and m base,
// the string range is: a-z, A-Z, 0-9 and'_','.' symbols
// see: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go?answertab=votes#tab-top
func RandNdigMbitString(n int, m ...int) string {
	b := make([]byte, n)
	dict := digits
	if 1 == len(m) {
		dict = digits[0:m[0]]
	} else if 2 == len(m) {
		dict = digits[m[0] : m[0]+m[1]]
	}

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

// New64BitID 获取 64进制 UUID, 10个数字+26个小写字母+26个大写字母+(=、_), 长度 24
// generates a random UUID according to RFC 4122
func New64BitID() string {
	return RandNdigMbitString(64)
}
