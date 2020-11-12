package main

import (
	"bytes"
	"errors"
	"text/template"
	"time"
)

// KEYS code 生成字典
const KEYS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

const (
	verificationTlp = `<body style="font-family: Roboto, sans-serif">
  <p>Hello, Welcome to whoam. You are using Email Verification Code to login to <a href="https://whoam.xyz">WHOAM</a>
  <p><big>Verification code: <b>{{ . }}</b>.</big>
  <p>It's valid within <b>15 minutes.</b>
  <p>If this isn't your own operating, please ignore this email.
  <p>Please don't reply!
    <hr>
  <p>Thank you,<p style="margin: 0 auto; font-size: 1.5em;">The ThreeTenth team
</body>`
)

const timeoutUserVerification = 900 // 用户验证码有效时长: 15分钟

// User basic information, id, email and
type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

// UserToken user login token
type UserToken struct {
	ID       uint   `gorm:"primarykey" schema:"-"`
	UserID   uint   `schema:"-"`
	DeivceID string `schema:"deivceId,required"`
	AppID    string `schema:"appId,required"`
	Token    string `schema:"-"`
	Uptoken  string `schema:"-"`
}

// UserVerification 用户登录验证信息
type UserVerification struct {
	Email     string `schema:"email,required"`
	Code      string `schema:"code,required"`
	Token     string `schema:"token,required"`
	UntilTime int64  `schema:"-"`
}

// Return verification failure information,
// if returns nil, it means verification is successful
func (user *UserVerification) verifica(dst *UserVerification) error {
	if user.UntilTime < time.Now().Unix() {
		return errors.New("Verification failed: code is expired")
	}
	if user.Code != dst.Code {
		return errors.New("Verification failed: code is invalid")
	}
	if user.Code != dst.Code {
		return errors.New("Verification failed: token is invalid")
	}
	if user.Email != dst.Email {
		return errors.New("Verification failed: email is invalid")
	}

	return nil
}

// 用户登录验证信息
var userVerificationMap map[string]UserVerification
var userTokenMap map[string]UserToken

func initUser() {
	db.AutoMigrate(&User{}, &UserToken{})
}

// PostUserAuth 用户登录授权验证
func PostUserAuth(c *Context) error {
	var user UserVerification
	err := c.ParseForm(&user)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	userVerification, ok := userVerificationMap[user.Token]

	if !ok {
		return c.Unauthorized("Verification failed: token is invalid")
	}

	err = userVerification.verifica(&user)

	if err != nil {
		return c.Unauthorized(err.Error())
	}

	// todo 记录用户验证成功

	return c.Ok("LOGIN SUCCESSED")
}

// PostUserLogin 用户登录
func PostUserLogin(c *Context) error {
	email, err := c.FormValue("email")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var userToken UserToken
	err = c.ParseForm(&userToken)

	code := genRandCode(4, KEYS[0:36])
	t, err := template.New("login").Parse(verificationTlp)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, code)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	err = SendMail(email, "Login WHOAM with verification code", buf.String())
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	token, _ := New64BitUUID()

	userVerificationMap[token] = UserVerification{email, code, token, time.Now().Unix() + timeoutUserVerification}
	userTokenMap[token] = userToken

	return c.Ok(token)
}
