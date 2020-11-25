package main

import (
	"errors"
	"math/rand"
	"net"
	"regexp"
	"time"
	"unsafe"
)

// ExternalIP 获取 IP
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

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

// New128BitID Get a 128-base random string,
// 10 numbers + 26 lowercase letters + 26 uppercase letters + (=, _), length 64
func New128BitID() string {
	return RandNdigMbitString(128)
}

// New64BitID Get a 64-base random string,
// 10 numbers + 26 lowercase letters + 26 uppercase letters + (=, _), length 64
func New64BitID() string {
	return RandNdigMbitString(64)
}

// New32BitID Get a 32-base random string,
// 10 numbers + 26 uppercase letters
func New32BitID() string {
	return RandNdigMbitString(32, 36)
}

// New32bitID Get a 32-base random string,
// 10 numbers + 26 lowercase letters
func New32bitID() string {
	return RandNdigMbitString(32, 26, 36)
}

// New16BitID Get a 16-base random string,
// 10 numbers + 26 uppercase letters
func New16BitID() string {
	return RandNdigMbitString(16, 36)
}

// New16bitID Get a 16-base random string,
// 10 numbers + 26 lowercase letters
func New16bitID() string {
	return RandNdigMbitString(16, 26, 36)
}

// New4BitID Get a 4-base random string,
// 10 numbers + 26 uppercase letters
func New4BitID() string {
	return RandNdigMbitString(4, 36)
}

// New4bitID Get a 4-base random string,
// 10 numbers + 26 lowercase letters
func New4bitID() string {
	return RandNdigMbitString(4, 26, 36)
}
