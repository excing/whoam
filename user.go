package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"
)

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
	ID          uint      `json:"-" gorm:"primarykey"`
	CreatedAt   time.Time `json:"-"`
	UserID      uint      `json:"userId"`
	AppID       string    `json:"-"`
	AccessToken string    `json:"accessToken"`
	UpdateToken string    `json:"updateToken"`
}

// Return verification failure information,
// if returns nil, it means verification is successful
func (user *userVerificationForm) verifica(dst *userVerificationForm) error {
	if user.UntilTime < time.Now().Unix() {
		return errors.New("Verification failed: code is expired")
	}
	if user.Code != strings.ToTitle(dst.Code) {
		return errors.New("Verification failed: code is invalid")
	}
	if user.Token != dst.Token {
		return errors.New("Verification failed: token is invalid")
	}
	if user.State != dst.State {
		return errors.New("Verification failed: state is invalid")
	}
	if user.Email != dst.Email {
		return errors.New("Verification failed: email is invalid")
	}

	return nil
}

// 用户登录验证信息
var userVerificationMap map[string]userVerificationForm
var oauthTokenMap map[string]UserToken

func initUser() {
	db.AutoMigrate(&User{}, &UserToken{})

	userVerificationMap = make(map[string]userVerificationForm)
	oauthTokenMap = make(map[string]UserToken)
}

type userVerificationForm struct {
	Email     string `schema:"email,required"`
	AppID     string `schema:"-"`
	State     string `schema:"state,required" note:"This parameter should be consistent with the state in /user/main/login"`
	Code      string `schema:"code,required"`
	Token     string `schema:"token,required"`
	UntilTime int64  `schema:"-"`
}

// PostUserAuth 用户登录授权验证
func PostUserAuth(c *Context) error {
	var dst userVerificationForm
	err := c.ParseForm(&dst)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	userVerification, ok := userVerificationMap[dst.Token]

	if !ok {
		return c.Unauthorized("Verification failed: token is invalid")
	}

	err = userVerification.verifica(&dst)

	if err != nil {
		return c.Unauthorized(err.Error())
	}

	// todo 记录用户验证成功
	var user User
	if db.Where("email=?", userVerification.Email).Find(&user).Error != nil || 0 == user.ID {
		user.Email = userVerification.Email
		if err = db.Create(&user).Error; err != nil {
			return c.InternalServerError(err.Error())
		}
	}

	accessToken := New64BitID()
	updateToken := New64BitID()

	userToken := UserToken{
		UserID:      user.ID,
		AppID:       userVerification.AppID,
		AccessToken: accessToken,
		UpdateToken: updateToken,
	}

	if err = db.Create(&userToken).Error; err != nil {
		return c.InternalServerError(err.Error())
	}

	delete(userVerificationMap, dst.Token)

	return c.Ok(&userToken)
}

type userLoginForm struct {
	Email string `schema:"email,required"`
	AppID string `schema:"appId,required"`
	State string `schema:"state,required" note:"random number"`
}

// PostMainCode 用户登录
func PostMainCode(c *Context) error {
	var form userLoginForm
	err := c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if !VerifyEmailFormat(form.Email) {
		return c.BadRequest("Email is invalid")
	}

	code := New4BitID()
	t, err := template.New("login").Parse(verificationTlp)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, code)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	// err = SendMail(email, "Login WHOAM with verification code", buf.String())
	fmt.Println(buf.String())
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	token := New64BitID()

	userVerificationMap[token] = userVerificationForm{form.Email, form.AppID, form.State, code, token, time.Now().Unix() + timeoutUserVerification}

	return c.Ok(token)
}

// GetUserState get user login status
func GetUserState(c *Context) error {
	var count int64
	token := c.GetHeader("Authorization")
	if err := db.Where("access_token=?", token).Find(&UserToken{}).Count(&count).Error; err != nil || 0 == count {
		return c.Any().Unauthorized("Invalid token, please login again")
	}
	return c.NoContent()
}

// UserOAuthLoginForm `/user/oauth/login` api form
type UserOAuthLoginForm struct {
}

// PostUserOAuthLogin requesting user's whoam identity
// GitHub: https://github.com/login?client_id=bfe378e98cde9624c98c&return_to=/login/oauth/authorize?client_id=bfe378e98cde9624c98c&redirect_uri=https://www.iconfont.cn/api/login/github/callback&state=123123sadh1as12
// 58778ef7632c0d4876432bb4866206775c711d4c
func PostUserOAuthLogin(c *Context) error {
	var form UserOAuthLoginForm
	err := c.ParseForm(&form)

	if err != nil {
		return c.BadRequest(err.Error())
	}

	return c.OkHTML(tlpUserOAuthLogin, form)
}

type oauthAuthForm struct {
	UserID      string `schema:"userId,required"`
	AccessToken string `schema:"accessToken,required"`
	ClientID    string `schema:"clientId,required"`
	State       string `schema:"state,required"`
}

// PostUserOAuthAuth whoam user authorized the request(/user/oauth/login request)
func PostUserOAuthAuth(c *Context) error {
	var form oauthAuthForm
	err := c.ParseForm(&form)

	if err != nil {
		return c.BadRequest(err.Error())
	}

	var loginUserToken UserToken
	if err = db.Where("user_id=? AND access_token=?", form.UserID, form.AccessToken).Find(&loginUserToken).Error; err != nil || 0 == loginUserToken.ID {
		return c.Any().Unauthorized("Invalid token, please login again")
	}

	accessToken := New64BitID()
	updateToken := New64BitID()

	oauthUserToken := UserToken{
		UserID:      loginUserToken.UserID,
		AppID:       form.ClientID,
		AccessToken: accessToken,
		UpdateToken: updateToken,
	}

	if err = db.Create(&oauthUserToken).Error; err != nil {
		return c.InternalServerError(err.Error())
	}

	code := New32bitID()

	oauthTokenMap[code] = oauthUserToken

	return c.Ok(code)
}

// GetOAuthCode obtain user authentication information through code
func GetOAuthCode(c *Context) error {
	code := c.Query("code")
	if "" == code {
		return c.BadRequest("code is empty")
	}

	userOAuthToken, ok := oauthTokenMap[code]

	if !ok {
		return c.Unauthorized("Invalid token, please login again")
	}

	delete(oauthTokenMap, code)

	return c.Ok(&userOAuthToken)
}

type oauthStateForm struct {
	UserID      string `schema:"userId,required"`
	AccessToken string `schema:"accessToken,required"`
	ClientID    string `schema:"clientId,required"`
}

// GetOAuthState Get user authorization status
func GetOAuthState(c *Context) error {
	var form oauthStateForm
	err := c.ParseForm(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var count int64
	if db.Where("user_id=? AND app_id=? AND access_token=?", form.UserID, form.ClientID, form.AccessToken).Find(&UserToken{}).Error != nil && 0 == count {
		return c.Unauthorized("Invalid token, please login again")
	}

	return c.NoContent()
}
