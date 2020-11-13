package main

import "regexp"

// VerifyEmailFormat Verify email address
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` // email regular expression
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
