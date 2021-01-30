package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"whoam.xyz/ent"
	"whoam.xyz/ent/service"
)

// loginEndpoint user loginEndpoint page
// Scenarios:
// 1. Once on the auth page, a `pushState` was performed and then refreshed again
// 2. Switch to the loginEndpoint page when the auth page is not logged in
// 3. Direct browser call
func loginEndpoint(c *Context) error {
	token := c.MustGet("token").(*StandardClaims)
	fmt.Println(token)
	return c.OkHTML(tlpUserLogin, nil)
}

// oauthEndpoint requesting user's whoam identity
// GitHub: https://github.com/login?client_id=bfe378e98cde9624c98c&return_to=/login/oauth/authorize?client_id=bfe378e98cde9624c98c&redirect_uri=https://www.iconfont.cn/api/login/github/callback&state=123123sadh1as12
// Scenarios:
// 1. Click from the A application  ✔
// 2. Refresh the auth page again
// 3. Click from the login page to come in
// 4. Invalid call
func oauthEndpoint(c *Context) error {
	var query struct {
		ClientID    string `form:"client_id" binding:"required"`
		RedirectURI string `form:"redirect_uri" binding:"required,url"`
		ReturnTo    string `form:"return_to"`
		State       string `form:"state" binding:"required"`
	}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_service, err := client.Service.Query().Where(service.IDEQ(query.ClientID)).First(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var response struct {
		User    *ent.User
		Service *ent.Service
	}

	token := c.MustGet("token").(*StandardClaims)
	if token == nil {
		return c.MovedPermanently("/user/login?" + c.Request.URL.RawQuery)
	}

	_user, err := client.User.Get(ctx, int(token.OtherID))
	if err != nil {
		return c.MovedPermanently("/user/login?" + c.Request.URL.RawQuery)
	}

	response.User = _user
	response.Service = _service

	return c.OkHTML(tlpUserOAuth, &response)
}

// AuthRequired middleware just in the "authorized" group.
func AuthRequired(c *gin.Context) {
	_jwtToken, err := FilterJWTToken(c.GetHeader("Authorization"), signingKey)
	if err != nil {
		if accessToken, err := c.Cookie("access_token"); err == nil {
			_jwtToken, _ = FilterJWTToken(accessToken, signingKey)
		}
	}

	if _jwtToken != nil {
		if _jwtToken.Audience != MainServiceID {
			_jwtToken = nil
		}
	}

	c.Set("token", _jwtToken)

	c.Next()
}
