package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"whoam.xyz/ent"
	"whoam.xyz/ent/method"
	"whoam.xyz/ent/service"
)

// PageUserLogin user login page
// Scenarios:
// 1. Once on the auth page, a `pushState` was performed and then refreshed again
// 2. Switch to the login page when the auth page is not logged in
// 3. Direct browser call
func PageUserLogin(c *Context) error {
	token := c.MustGet("token").(*StandardClaims)
	fmt.Println(token)
	return c.OkHTML(tlpUserLogin, nil)
}

// PageUserOAuth requesting user's whoam identity
// GitHub: https://github.com/login?client_id=bfe378e98cde9624c98c&return_to=/login/oauth/authorize?client_id=bfe378e98cde9624c98c&redirect_uri=https://www.iconfont.cn/api/login/github/callback&state=123123sadh1as12
// Scenarios:
// 1. Click from the A application  ✔
// 2. Refresh the auth page again
// 3. Click from the login page to come in
// 4. Invalid call
func PageUserOAuth(c *Context) error {
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

	token := c.MustGet("token").(*StandardClaims)
	if token == nil {
		return c.OkHTML(tlpUserOAuth, nil)
	}

	_service, err := client.Service.Query().Where(service.IDEQ(query.ClientID)).First(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_permissions, err := _service.QueryPermissions().Where(method.ScopeEQ(method.ScopePrivate)).All(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_result := make(map[string][]*ent.Method)
	for _, _permission := range _permissions {
		_service, err := _permission.QueryOwner().First(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		_val, ok := _result[_service.Name]
		_val = append(_val, _permission)
		if !ok {
			_result[_service.Name] = _val
		}
	}

	return c.OkHTML(tlpUserOAuth, &_result)
}

func authorizeUser(c *gin.Context) {
	_claims, err := FilterJWTToken(c.GetHeader("Authorization"), signingKey)
	if err != nil {
		if accessToken, err := c.Cookie("accessToken"); err == nil {
			_claims, _ = FilterJWTToken(accessToken, signingKey)
		}
	}

	c.Set("token", _claims)

	c.Next()
}
