package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 授权状态：等待中
const statusAuthorizationWaiting = 0

// 授权状态：同意授权
const statusAuthorizationUserGrant = 1

// 授权状态：拒绝授权
const statusAuthorizationUserDeny = -1

const timeoutTokenExpire = 604800 // 授权有效期为 7 天
const timeoutTokenDelete = 259200 // 授权失效后，保留 UserToken 3 天

// ServiceInfo 服务信息
type ServiceInfo struct {
	// 服务 ID
	ServiceID string
	// 服务名
	ServiceName string
	// 服务详细，可选
	ServiceDesc string
	// 服务 Token 的 crypto/bcrypt 值
	ServiceTokenEncode string
}

// AuthorizationInfo 认证请求信息
// 请求用户将 ServiceProvider 服务授权给 ServiceRequester
type AuthorizationInfo struct {
	// 用户名
	UserName string
	// 用户 Token 的 crypto/bcrypt 值
	UserTokenEncode string
	// 提供方服务ID
	ServiceProviderID string
	// 请求方服务ID
	ServiceRequesterID string
	// 服务授权状态
	// see: STATUS_AUTHORIZATION_WAITING, STATUS_AUTHORIZATION_USER_GRANT, STATUS_AUTHORIZATION_USER_DENY
	AuthorizationStatus int

	// 用户 Token 创建时间
	TokenCreateTime int64
	// 用户 Token 刷新时间
	TokenUpdateTime int64
	// 用户 Token 过期时间
	TokenExpireTime int64
	// 用户 Token 将被删除的时间
	TokenDeleteTime int64
}

// 授权请求列表. key 为 authorizationCode, value 为用户 Token 的用户信息, 包含授权状态
var authorizationMap map[string]AuthorizationInfo

// 服务提供者列表
var serviceMap map[string]ServiceInfo

// 请求授权邮件列表
var authorizationEmailMap map[string]string

// PostServicer 提交服务注册
func PostServicer(c *Context) error {
	serviceID := c.PostForm("serviceId")
	serviceName := c.PostForm("serviceName")
	serviceDesc := c.PostForm("serviceDesc")

	if "" == serviceID {
		return c.BadRequest("serviceId is empty")
	} else if "" == serviceName {
		return c.BadRequest("serviceName is empty")
	}

	var providerRegisterResult, requesterRegisterResult = false, false

	registerResult := make(map[string]string)

	if _, ok := serviceMap[serviceID]; !ok {
		token, encodeToken, err := "", "", errors.New("")

		if err != nil {
			return c.InternalServerError("registration failed: %s", err.Error())
		}

		serviceMap[serviceID] = ServiceInfo{
			serviceID,
			serviceName,
			serviceDesc,
			encodeToken,
		}
		registerResult["providerToken"] = token

		providerRegisterResult = true
	}

	if providerRegisterResult || requesterRegisterResult {
		return c.Ok(registerResult)
	}
	return c.Created("This service is already registered")
}

// DeleteServicer 注销服务
func DeleteServicer(c *Context) error {
	serviceID := c.PostForm("serviceId")
	serviceToken := c.PostForm("serviceToken")

	if "" == serviceID {
		return c.BadRequest("serviceId is empty")
	} else if "" == serviceToken {
		return c.BadRequest("serviceToken is empty")
	}

	serviceProvider, providerOk := serviceMap[serviceID]

	if providerOk {
		err := bcrypt.CompareHashAndPassword([]byte(serviceProvider.ServiceTokenEncode), []byte(serviceToken))

		if err == nil {
			delete(serviceMap, serviceID)
			return c.NoContent()
		}
	}

	return c.Unauthorized("serviceToken is an invalid value")
}
