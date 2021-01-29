package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"whoam.xyz/ent/oauth"
	"whoam.xyz/ent/user"
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

const timeoutUserVerification = 900             // 用户验证码有效时长: 15分钟
const timeoutRefreshToken = 30 * 24 * time.Hour // user refresh token timeout: 30day
const timeoutAccessToken = 7 * time.Minute      // user access token timeout: 7min

// 用户登录验证信息
var userVerificaBox *Box
var oauthCodeBox *Box
var signingKey []byte

// InitUser initialize User related
func InitUser() {
	// size: 3M
	// default timeout: 15min
	userVerificaBox = NewBox(3*1024*1024, 15*60)
	// size: 3M
	// default timeout: 5min
	oauthCodeBox = NewBox(3*1024*1024, 5*60)

	signingKey = []byte(New32BitID())
}

// UserAuthorize Return true, if the specified accessToken is not found, return false
func UserAuthorize(accessToken string) bool {
	_, err := FilterJWTToken(accessToken, signingKey)
	if err != nil {
		return false
	}

	return true
}

type userVerificationForm struct {
	Email string `json:"email" binding:"required"`
	State string `json:"state" binding:"required" note:"This parameter should be consistent with the state in /user/main/login"`
	Code  string `json:"code" binding:"required"`
	Token string `json:"token" binding:"required"`
}

// PostUserAuth 用户登录授权验证
func PostUserAuth(c *Context) error {
	var dst userVerificationForm
	err := c.ShouldBindJSON(&dst)
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

	user, err := client.User.Query().Where(user.EmailEQ(src.Email)).Only(ctx)
	if err != nil {
		user, err = client.User.Create().SetEmail(src.Email).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}
	}

	// accessToken := New64BitID()
	accessToken, err := NewJWTToken(user.ID, MainServiceID, timeoutAccessToken, signingKey)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	mainToken := New64BitID()

	auth, err := client.Oauth.Create().
		SetMainToken(mainToken).
		SetExpiredAt(time.Now().Add(timeoutRefreshToken)).
		SetOwner(user).
		SetServiceID(MainServiceID).
		Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	userVerificaBox.DelString(dst.Token)

	return c.Ok(
		struct {
			AccessToken string `json:"accessToken"`
			MainToken   string `json:"mainToken"`
		}{
			AccessToken: accessToken,
			MainToken:   auth.MainToken,
		})
}

type userLoginForm struct {
	Email string `json:"email" binding:"required"`
	State string `json:"state" binding:"required" note:"random number"`
}

// PostMainCode 用户登录
func PostMainCode(c *Context) error {
	var form userLoginForm
	err := c.ShouldBindJSON(&form)
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

// GetUser is used to get the basic information of the user, if authenticated through OAuth
func GetUser(c *Context) error {
	accessToken, _ := c.Cookie("accessToken")
	if "" == accessToken {
		accessToken = c.GetHeader("Authorization")
		if "" == accessToken {
			return c.Unauthorized("Unauthorized")
		}
	} else if accessToken != c.GetHeader("Authorization") {
		// 冲突
		return c.Conflict("Cookie's accessToken and Header's Authorization value are inconsistent")
	}

	_claims, err := FilterJWTToken(accessToken, signingKey)
	if err != nil {
		return c.Unauthorized(err.Error())
	}

	_user, err := client.User.Query().Where(user.IDEQ(int(_claims.OtherID))).Only(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(
		struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}{
			ID:    _user.ID,
			Email: _user.Email,
		})
}

// PostUserOAuthAuth whoam user authorized the request(/user/oauth/auth request)
func PostUserOAuthAuth(c *Context) error {
	var form struct {
		MainToken   string `json:"mainToken" binding:"required"`
		ClientID    string `json:"clientId" binding:"required"`
		State       string `json:"state" binding:"required"`
		Permissions []int  `json:"permissions" binding:"required"`
	}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	owner, err := client.Oauth.Query().Where(oauth.MainTokenEQ(form.MainToken)).QueryOwner().Only(ctx)
	if err != nil {
		return c.Unauthorized("Invalid token, please login again")
	}

	oauthUser := userOAuth{
		UserID:   owner.ID,
		ClientID: form.ClientID,
	}

	code := New32bitID()
	if err = oauthCodeBox.SetVal(code, &oauthUser); err != nil {
		return c.InternalServerError(err.Error())
	}

	err = oauthCodeBox.Val(code, &oauthUser)

	return c.Ok(code)
}

type userOAuth struct {
	UserID   int    `json:"userId"`
	ClientID string `json:"clientId"`
}

// GetOAuthCode obtain user authentication information through code
func GetOAuthCode(c *Context) error {
	code := c.Query("code")
	if "" == code {
		return c.BadRequest("code is empty")
	}

	var oauthUser userOAuth

	err := oauthCodeBox.Val(code, &oauthUser)
	if err != nil {
		return c.Unauthorized("Invalid token, please login again")
	}

	oauthCodeBox.DelString(code)

	accessToken, err := NewJWTToken(oauthUser.UserID, oauthUser.ClientID, timeoutAccessToken, signingKey)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	auth, err := client.Oauth.Create().
		SetMainToken(New64BitID()).
		SetExpiredAt(time.Now().Add(timeoutRefreshToken)).
		SetOwnerID(oauthUser.UserID).
		SetServiceID(oauthUser.ClientID).
		Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(
		struct {
			AccessToken string `json:"accessToken"`
			MainToken   string `json:"mainToken"`
		}{
			AccessToken: accessToken,
			MainToken:   auth.MainToken,
		})
}

// GetOAuthState Get user authorization status
func GetOAuthState(c *Context) error {
	if !UserAuthorize(c.GetHeader("Authorization")) {
		return c.Unauthorized("Invalid token, please login again")
	}

	return c.NoContent()
}

type oauthRefreshForm struct {
	UserID       uint   `schema:"userId" binding:"required"`
	RefreshToken string `schema:"refreshToken" binding:"required"`
	ClientID     string `schema:"clientId" binding:"required"`
}

// PostUserOAuthRefresh refresh user access token
func PostUserOAuthRefresh(c *Context) error {
	refreshToken, err := c.GetFormString("refreshToken")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	auth, err := client.Oauth.Query().
		Where(oauth.MainTokenEQ(refreshToken)).
		Where(oauth.ExpiredAtGT(time.Now())).
		Only(ctx)
	if err != nil {
		return c.Unauthorized("Invalid refreshToken, please login again")
	}

	authUser, err := auth.QueryOwner().Only(ctx)
	if err != nil {
		return c.Unauthorized("Invalid authorized user, please login again")
	}

	authService, err := auth.QueryService().Only(ctx)
	if err != nil {
		return c.Unauthorized("Invalid authorized service, please login again")
	}

	accessToken, err := NewJWTToken(authUser.ID, authService.ID, timeoutAccessToken, signingKey)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	_, err = auth.Update().SetExpiredAt(time.Now().Add(timeoutRefreshToken)).Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(accessToken)
}
