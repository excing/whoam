package main

import (
	"bytes"
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

const timeoutUserVerification = 900           // 用户验证码有效时长: 15分钟
const timeoutUserToken = 7 * 24 * time.Second // user token timeout: 7 day

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
	ExpiredAt   int64     `json:"-"`
	UserID      uint      `json:"userId"`
	ServiceID   string    `json:"serviceId"`
	AccessToken string    `json:"accessToken"`
	UpdateToken string    `json:"updateToken"`
}

// 用户登录验证信息
var userVerificaBox *Box
var oauthCodeBox *Box

func initUser() {
	db.AutoMigrate(&User{}, &UserToken{})

	// size: 3M
	// default timeout: 15min
	userVerificaBox = NewBox(3*1024*1024, 15*60)
	// size: 3M
	// default timeout: 5min
	oauthCodeBox = NewBox(3*1024*1024, 5*60)
}

type userVerificationForm struct {
	Email string `schema:"email,required"`
	State string `schema:"state,required" note:"This parameter should be consistent with the state in /user/main/login"`
	Code  string `schema:"code,required"`
	Token string `schema:"token,required"`
}

// PostUserAuth 用户登录授权验证
func PostUserAuth(c *Context) error {
	var dst userVerificationForm
	err := c.ParseForm(&dst)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var src userVerificationForm
	err = userVerificaBox.Val(dst.Token, &src)

	if err != nil {
		return c.Unauthorized("Verification failed: token is invalid or code is expired")
	}

	if src.Code != strings.ToTitle(dst.Code) {
		return c.Unauthorized("Verification failed: code is invalid")
	}
	if src.Token != dst.Token {
		return c.Unauthorized("Verification failed: token is invalid")
	}
	if src.State != dst.State {
		return c.Unauthorized("Verification failed: state is invalid")
	}
	if src.Email != dst.Email {
		return c.Unauthorized("Verification failed: email is invalid")
	}

	var user User
	if db.Where("email=?", src.Email).Find(&user).Error != nil || 0 == user.ID {
		user.Email = src.Email
		if err = db.Create(&user).Error; err != nil {
			return c.InternalServerError(err.Error())
		}
	}

	accessToken := New64BitID()
	updateToken := New64BitID()

	userToken := UserToken{
		ExpiredAt:   time.Now().Add(timeoutUserToken).UnixNano(),
		UserID:      user.ID,
		ServiceID:   MainServiceID,
		AccessToken: accessToken,
		UpdateToken: updateToken,
	}

	if err = db.Create(&userToken).Error; err != nil {
		return c.InternalServerError(err.Error())
	}

	userVerificaBox.DelString(dst.Token)

	return c.Ok(&userToken)
}

type userLoginForm struct {
	Email string `schema:"email,required"`
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
	err = userVerificaBox.SetVal(token, userVerificationForm{form.Email, form.State, code, token})
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(token)
}

// UserOAuthLoginForm `/user/oauth/login` api form
type UserOAuthLoginForm struct {
}

// PageUserLogin user login page
func PageUserLogin(c *Context) error {
	return c.OkHTML(tlpUserLogin, nil)
}

// PageUserOAuth requesting user's whoam identity
// GitHub: https://github.com/login?client_id=bfe378e98cde9624c98c&return_to=/login/oauth/authorize?client_id=bfe378e98cde9624c98c&redirect_uri=https://www.iconfont.cn/api/login/github/callback&state=123123sadh1as12
func PageUserOAuth(c *Context) error {
	if _, ok := c.GetQuery("clientId"); !ok {
		return c.BadRequest("clientId is empty")
	}

	c.Writer.Header().Set("Cache-control", "private")

	if accessToken, err := c.Cookie("accessToken"); err != nil || !UserAuthorized(MainServiceID, accessToken) {
		returnTo, ok := c.GetQuery("return_to")
		url := c.Request.URL
		if !ok {
			returnTo = c.GetHeader("Referer")
			return c.MovedPermanently("/user/login?" + url.RawQuery + "&return_to=" + returnTo)
		}
		return c.MovedPermanently("/user/login?" + url.RawQuery)
	}

	return c.OkHTML(tlpUserOAuth, nil)
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
		return c.Unauthorized("Invalid token, please login again")
	}

	accessToken := New64BitID()
	updateToken := New64BitID()

	oauthUserToken := UserToken{
		UserID:      loginUserToken.UserID,
		ServiceID:   form.ClientID,
		AccessToken: accessToken,
		UpdateToken: updateToken,
	}

	code := New32bitID()
	if err = oauthCodeBox.SetVal(code, &oauthUserToken); err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(code)
}

// GetOAuthCode obtain user authentication information through code
func GetOAuthCode(c *Context) error {
	code := c.Query("code")
	if "" == code {
		return c.BadRequest("code is empty")
	}

	var oauthUserToken UserToken

	err := oauthCodeBox.Val(code, &oauthUserToken)
	if err != nil {
		return c.Unauthorized("Invalid token, please login again")
	}

	oauthCodeBox.DelString(code)

	oauthUserToken.ExpiredAt = time.Now().Add(timeoutUserToken).UnixNano()

	if err = db.Create(&oauthUserToken).Error; err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(&oauthUserToken)
}

type oauthStateForm struct {
	UserID      string `schema:"userId,required"`
	AccessToken string `schema:"accessToken,required"`
	ClientID    string `schema:"clientId,required"`
}

// GetOAuthState Get user authorization status
func GetOAuthState(c *Context) error {
	accessToken := c.GetHeader("Authorization")
	clientID := c.Query("client_id")

	var user UserToken
	if db.Where("service_id=? AND access_token=? AND ?<expired_at", clientID, accessToken, time.Now().UnixNano()).Find(&user).Error != nil || 0 == user.ID {
		return c.Unauthorized("Invalid token, please login again")
	}

	return c.NoContent()
}

// UserAuthorized Return true, if the specified clientID and accessToken are not found, return false
func UserAuthorized(clientID string, accessToken string) bool {
	var user UserToken
	if db.Where("service_id=? AND access_token=?", clientID, accessToken).Find(&user).Error != nil || 0 == user.ID {
		return false
	}
	return true
}
