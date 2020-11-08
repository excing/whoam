package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

// KEYS code 生成字典
const KEYS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

var decoder = schema.NewDecoder()

func inout(handle func(p *Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := &Context{c}

		proxy = proxy.Any()

		handle(proxy)
	}
}

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
		token, encodeToken, err := genServiceToken(serviceID)

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

// PostAuthRequest 提交一个授权请求
func PostAuthRequest(c *Context) error {

	userName := c.PostForm("username")
	serviceProviderID := c.PostForm("providerId")
	serviceRequesterID := c.PostForm("requesterId")

	if "" == userName {
		return c.BadRequest("username is empty")
	} else if "" == serviceProviderID {
		return c.BadRequest("providerId is empty")
	} else if "" == serviceRequesterID {
		return c.BadRequest("requesterId is empty")
	} else if serviceProviderID == serviceRequesterID {
		return c.BadRequest("providerId equals requesterId")
	}

	serviceProvider, serviceProviderIDOk := serviceMap[serviceProviderID]
	serviceRequester, serviceRequesterIDOk := serviceMap[serviceRequesterID]

	if !serviceProviderIDOk {
		return c.Forbidden("providerId does not exist")
	} else if !serviceRequesterIDOk {
		return c.Forbidden("requesterId does not exist")
	}

	authorizationCode, userToken, userTokenEncode, err := genAuthCodeAndToken(userName, serviceProvider, serviceRequester)

	if err != nil {
		return c.InternalServerError("Authorization failed: %s", err.Error())
	}

	timestamp := time.Now().Unix()

	authorizationMap[authorizationCode] = AuthorizationInfo{
		userName,
		userTokenEncode,
		serviceProviderID,
		serviceRequesterID,
		statusAuthorizationWaiting,
		timestamp,
		timestamp,
		timestamp + timeoutTokenExpire,
		timestamp + timeoutTokenExpire + timeoutTokenDelete,
	}

	go sendAuthorizationEmail(authorizationCode, authorizationMap[authorizationCode], userToken)

	data := make(map[string]string)
	data["authCode"] = authorizationCode
	data["userToken"] = userToken

	return c.Ok(data)
}

// GrantAuthRequest 同意一个授权请求
func GrantAuthRequest(c *Context) error {
	code, _ := c.Params.Get("code")

	if "" == code {
		return c.BadRequest("code is empty")
	}

	authorizationCode, ok := authorizationEmailMap[code]

	if !ok {
		return c.NotFound("code is an invalid value")
	}

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		return c.NotFound("code is an invalid value")
	}

	fmt.Println("Grant Authorization", code, authorizationCode)

	authorizationInfo.AuthorizationStatus = statusAuthorizationUserGrant

	authorizationMap[authorizationCode] = authorizationInfo

	updateUserTokenTime(authorizationInfo)

	return c.NoContent()
}

// DenyAuthRequest 拒绝一个授权请求
func DenyAuthRequest(c *Context) error {
	code, _ := c.Params.Get("code")

	if "" == code {
		return c.BadRequest("code is empty")
	}

	authorizationCode, ok := authorizationEmailMap[code]

	if !ok {
		return c.NotFound("code is an invalid value")
	}

	delete(authorizationEmailMap, code)

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		return c.NotFound("code is an invalid value")
	}

	fmt.Println("Deny Authorization", code, authorizationCode)

	authorizationInfo.AuthorizationStatus = statusAuthorizationUserDeny

	authorizationMap[authorizationCode] = authorizationInfo

	updateUserTokenTime(authorizationInfo)

	return c.NoContent()
}

// UpgradeAuthRequest 刷新一个授权令牌
func UpgradeAuthRequest(c *Context) error {
	authorizationCode := c.PostForm("authCode")
	userToken := c.PostForm("userToken")

	if "" == authorizationCode {
		return c.BadRequest("authCode is empty")
	} else if "" == userToken {
		return c.BadRequest("userToken is empty")
	}

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		return c.NotFound("authCode is an invalid value")
	}

	err := bcrypt.CompareHashAndPassword([]byte(authorizationInfo.UserTokenEncode), []byte(userToken))

	if err != nil {
		return c.Unauthorized("userToken is an invalid value")
	}

	timestamp := time.Now().Unix()

	if authorizationInfo.TokenExpireTime < timestamp {
		if authorizationInfo.TokenDeleteTime < timestamp {
			delete(authorizationMap, authorizationCode)
		}
		return c.Unauthorized("userToken has expired")
	}

	authorizationStatus := authorizationInfo.AuthorizationStatus

	if statusAuthorizationUserDeny == authorizationStatus {
		return c.Unauthorized("The user denied access")
	}

	updateUserTokenTime(authorizationInfo)

	return c.NoContent()
}

// GetAuthState 获取授权状态
func GetAuthState(c *Context) error {
	authorizationCode := c.PostForm("authCode")
	userToken := c.PostForm("userToken")

	if "" == authorizationCode {
		return c.BadRequest("authCode is empty")
	} else if "" == userToken {
		return c.BadRequest("userToken is empty")
	}

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		return c.NotFound("authCode is an invalid value")
	}

	err := bcrypt.CompareHashAndPassword([]byte(authorizationInfo.UserTokenEncode), []byte(userToken))

	if err != nil {
		return c.Unauthorized("userToken is an invalid value")
	}

	timestamp := time.Now().Unix()

	t1 := strconv.FormatInt(timestamp-authorizationInfo.TokenExpireTime, 10)
	t2 := strconv.FormatInt(timestamp-authorizationInfo.TokenDeleteTime, 10)

	fmt.Println(t1, t2)

	if authorizationInfo.TokenExpireTime < timestamp {
		if authorizationInfo.TokenDeleteTime < timestamp {
			delete(authorizationMap, authorizationCode)
		}
		return c.Unauthorized("userToken has expired")
	}

	authorizationStatus := authorizationInfo.AuthorizationStatus

	return c.Ok(map[string]int{"authStatus": authorizationStatus})
}

func updateUserTokenTime(authorizationInfo AuthorizationInfo) {
	timestamp := time.Now().Unix()
	authorizationInfo.TokenUpdateTime = timestamp
	authorizationInfo.TokenExpireTime = timestamp + timeoutTokenExpire
	authorizationInfo.TokenDeleteTime = timestamp + timeoutTokenExpire + timeoutTokenDelete
}

func genAuthCodeAndToken(userName string, serviceProvider ServiceInfo, serviceRequester ServiceInfo) (string, string, string, error) {
	UUID := genUUID()

	h := sha256.New()
	h.Write([]byte(UUID.String() + "-" + userName + "-" + serviceProvider.ServiceID + "-" + serviceRequester.ServiceID))
	token := hex.EncodeToString(h.Sum(nil))
	encodeToken, berr := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	if berr != nil {
		return "", "", "", errors.New("generate token has failure: " + userName)
	}
	return genAuthCode(), token, string(encodeToken), nil
}

func genServiceToken(serviceID string) (string, string, error) {
	UUID := genUUID().String()

	h := sha256.New()
	h.Write([]byte(UUID + "-" + serviceID))
	token := hex.EncodeToString(h.Sum(nil))
	encodeToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	if err != nil {
		return "", "", errors.New("generate token has failure: " + serviceID)
	}
	return token, string(encodeToken), nil
}

func genAuthCode() string {
	code := genRandCode(8, KEYS)

	if _, ok := authorizationMap[code]; ok {
		return genAuthCode()
	}

	return code
}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var r = rand.NewSource(time.Now().UnixNano())

// see: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go?answertab=votes#tab-top
func genRandCode(n int, dict string) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(dict) {
			b[i] = dict[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func genUUID() uuid.UUID {
	UUID, uerr := uuid.NewRandom()

	if uerr != nil {
		UUID = uuid.New()
	}

	return UUID
}

func sendAuthorizationEmail(authorizationCode string, authorizationInfo AuthorizationInfo, userToken string) {
	emailCode := genUUID().String()

	grantURL := "http://" + serverDomain + "/v1/auth/grant/" + emailCode
	denyURL := "http://" + serverDomain + "/v1/auth/deny/" + emailCode

	serviceProvider, serviceProviderIDOk := serviceMap[authorizationInfo.ServiceProviderID]
	serviceRequester, serviceRequesterIDOk := serviceMap[authorizationInfo.ServiceRequesterID]

	if !serviceProviderIDOk || !serviceRequesterIDOk {
		fmt.Println("Send authorization email failure and emailCode is", authorizationCode)
		return
	}

	to := authorizationInfo.UserName

	toSplit := strings.Split(to, "@")
	if 2 != len(toSplit) {
		fmt.Println("UserName is an invalid value")
		return
	}

	fristAndSecondName := strings.Split(toSplit[0], ".")
	callName := strings.Join(fristAndSecondName, " ")
	callTime := time.Now().Format("Mon Jan _2 15:04:05 2006")

	subject := "允许" + serviceRequester.ServiceName + "访问您的" + serviceProvider.ServiceName + "服务吗"
	body := "hi, " + callName +
		":<p>允许" +
		serviceRequester.ServiceName +
		"(" + serviceRequester.ServiceID + ")" +
		"访问您的" + serviceProvider.ServiceName +
		"(" + serviceProvider.ServiceID + ")" + "服务吗<p>" +
		"如果允许访问，请点击<p>" +
		"<a href=\"" + grantURL + "\">" + "允许</a><p>" +
		"如果不允许访问，请点击<p>" +
		"<a href=\"" + denyURL + "\">" + "拒绝</a><p>" +
		"允许之后可以选择拒绝，拒绝之后无法选择允许，请周知。<p>请勿回复本邮件，谢谢<p>" +
		"<div style=\"text-align: right\">whoam<p>Asia/Shanghai " + callTime + "</p></div>"

	fmt.Println("Do you allow", serviceRequester.ServiceID, "to access", serviceProvider.ServiceID)

	err := SendMail(to, subject, body)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("to", to, "mail send finished")
		authorizationEmailMap[emailCode] = authorizationCode
	}
}

var serverPort int
var serverDomain string

func init() {
	ip, err := ExternalIP()
	if err != nil {
		fmt.Println(err)
	}

	port := 8030

	flag.IntVar(&serverPort, "p", port, "Authorization server port")
	flag.StringVar(&serverDomain, "h", ip+":"+strconv.Itoa(port), "Authorization server domain")
}

func main() {
	flag.Parse()
	flag.Usage()
	time.FixedZone("CST", 8*3600)

	if serverDomain == "" {
		fmt.Println("ServerDomain is empty")
		os.Exit(-1)
	}

	fmt.Println("ServerPort: ", serverPort)

	// SetupMailCredentials("Enter e-mail username: ", "Enter e-mail password: ")

	fmt.Println("Whoam is working")

	authorizationMap = make(map[string]AuthorizationInfo)
	serviceMap = make(map[string]ServiceInfo)
	authorizationEmailMap = make(map[string]string)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {})

	v1 := router.Group("/v1")
	v1.POST("/user/login", inout(PostUserLogin))
	v1.POST("/user/auth", inout(PostUserAuth))

	v1.POST("/servicer", inout(PostServicer))
	v1.DELETE("/servicer", inout(DeleteServicer))

	v1.POST("/auth/request", inout(PostAuthRequest))
	v1.GET("/auth/grant/:code", inout(GrantAuthRequest))
	v1.GET("/auth/deny/:code", inout(DenyAuthRequest))
	v1.POST("/auth/upgrade", inout(UpgradeAuthRequest))
	v1.GET("/auth/state", inout(GetAuthState))

	router.Run(":" + strconv.Itoa(serverPort))
}

// PostUserAuth 用户登录授权验证
func PostUserAuth(c *Context) error {
	err := c.Request.ParseForm()
	if err != nil {
		return c.BadRequest(err.Error())
	}

	type User struct {
		Email string `schema:"email,required"`
		Token string `schema:"token"`
		Code  string
	}

	var user User
	err = decoder.Decode(&user, c.Request.PostForm)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	return c.Ok(user)
}

// PostUserLogin 用户登录
func PostUserLogin(c *Context) error {
	email := c.PostForm("email")

	if "" == email {
		return c.BadRequest("email is empty")
	}

	return c.Ok(email)
}
