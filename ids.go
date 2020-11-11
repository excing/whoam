package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_="

var rander = crand.Reader // random function

// Encry 对 bytes 进行 sha256 编码，并对 sha256 编码进行随机盐加密
func Encry(bytes []byte) ([]byte, []byte, error) {
	h := sha256.New()
	h.Write(bytes)
	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	psd, err := bcrypt.GenerateFromPassword(dst, bcrypt.DefaultCost)
	return dst, psd, err
}

func genAuthCode() string {
	code := genRandCode(8, KEYS)

	if _, ok := authorizationMap[code]; ok {
		return genAuthCode()
	}

	return code
}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var r = rand.NewSource(time.Now().UnixNano())

// see: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go?answertab=votes#tab-top
func genRandCode(n int, dict string) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(dict) {
			b[i] = dict[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func genUUID() uuid.UUID {
	UUID, uerr := uuid.NewRandom()

	if uerr != nil {
		UUID = uuid.New()
	}

	return UUID
}

// New64BitUUID 获取 64进制 UUID, 10个数字+26个小写字母+26个大写字母+(=、_), 长度 24
// generates a random UUID according to RFC 4122
func New64BitUUID() (string, error) {
	uuid := make([]byte, 18)
	n, err := io.ReadFull(rander, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuidValueOf(uuid), nil
}

func uuidValueOf(b []byte) string {
	carry := 6
	size := 3

	var temp []byte
	var s [68]byte
	sindex := 0
	var start, end int

	for {
		if start == len(b) {
			break
		}
		end = start + size
		if len(b) < end {
			end = len(b)
		}
		temp = b[start:end]
		start = end
		u := uint64(temp[len(temp)-1])
		ucarry := 8
		for i := len(temp) - 2; 0 <= i; i-- {
			u |= uint64(temp[i]) << ucarry
			ucarry += 8
		}
		i := 4
		for {
			if u <= 0 {
				break
			}

			i--
			s[sindex+i] = digits[u&63]
			u >>= carry
		}
		for 0 < i {
			i--
			s[sindex+i] = digits[0]
		}
		sindex += 4
	}
	return string(s[0:sindex])
}
