package main

import (
	"time"

	"whoam.xyz/ent/oauth"
	"whoam.xyz/ent/user"
)

type userOAuth struct {
	UserID   int    `json:"userId"`
	ClientID string `json:"clientId"`
}

// GetOAuthState Get user authorization status
func GetOAuthState(c *Context) error {
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

	_, err := FilterJWTToken(accessToken, signingKey)
	if err != nil {
		return c.Unauthorized(err.Error())
	}

	return c.NoContent()
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
		MainToken string `json:"mainToken" binding:"required"`
		ClientID  string `json:"clientId" binding:"required"`
		State     string `json:"state" binding:"required"`
	}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	owner, err := client.Oauth.Query().Where(oauth.MainTokenEQ(form.MainToken)).QueryUser().Only(ctx)
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
		SetUserID(oauthUser.UserID).
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

// PostUserOAuthRefresh refresh user access token
func PostUserOAuthRefresh(c *Context) error {
	var _body struct {
		MainToken string `json:"mainToken" binding:"required"`
	}
	err := c.ShouldBindJSON(&_body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	auth, err := client.Oauth.Query().
		Where(oauth.MainTokenEQ(_body.MainToken)).
		Where(oauth.ExpiredAtGT(time.Now())).
		Only(ctx)
	if err != nil {
		return c.Unauthorized("Invalid refreshToken, please login again")
	}

	authUser, err := auth.QueryUser().Only(ctx)
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
