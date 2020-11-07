package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
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

// Result 返回消息结构
type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// 授权请求列表. key 为 authorizationCode, value 为用户 Token 的用户信息, 包含授权状态
var authorizationMap map[string]AuthorizationInfo

// 服务提供者列表
var serviceMap map[string]ServiceInfo

// 请求授权邮件列表
var authorizationEmailMap map[string]string

// AuthorizationRequestHandler /authorization/request?username=&providerId=&requesterId=
func AuthorizationRequestHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.PostFormValue("username")
	serviceProviderID := r.PostFormValue("providerId")
	serviceRequesterID := r.PostFormValue("requesterId")

	if "" == userName {
		wirteError(w, -1, errors.New("username is empty"))
		return
	} else if "" == serviceProviderID {
		wirteError(w, -2, errors.New("providerId is empty"))
		return
	} else if "" == serviceRequesterID {
		wirteError(w, -3, errors.New("requesterId is empty"))
		return
	} else if serviceProviderID == serviceRequesterID {
		wirteError(w, -6, errors.New("providerId equals requesterId"))
		return
	}

	serviceProvider, serviceProviderIDOk := serviceMap[serviceProviderID]
	serviceRequester, serviceRequesterIDOk := serviceMap[serviceRequesterID]

	if !serviceProviderIDOk {
		wirteError(w, -4, errors.New("Can not find the service pointed to by providerId"))
		return
	} else if !serviceRequesterIDOk {
		wirteError(w, -5, errors.New("Can not find the service pointed to by requesterId"))
		return
	}

	authorizationCode, userToken, userTokenEncode, err := genAuthCodeAndToken(userName, serviceProvider, serviceRequester)

	if err != nil {
		fmt.Println("Request authorization failed: ", userName, serviceProviderID, serviceRequesterID)
		wirteError(w, -7, errors.New("Request authorization failed"))
		return
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

	wirteBody(w, 1, data)
}

// AuthorizationGrantHandler 允许授权
func AuthorizationGrantHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	if "" == code {
		wirteError(w, -1, errors.New("code is empty"))
		return
	}

	authorizationCode, ok := authorizationEmailMap[code]

	if !ok {
		wirteError(w, -2, errors.New("code is an invalid value"))
		return
	}

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		wirteError(w, -3, errors.New("code is an invalid value"))
		return
	}

	fmt.Println("Grant Authorization", code, authorizationCode)

	authorizationInfo.AuthorizationStatus = statusAuthorizationUserGrant

	authorizationMap[authorizationCode] = authorizationInfo

	updateUserTokenTime(authorizationInfo)

	wirteBody(w, 1, "Authorized success")
}

// AuthorizationDenyHandler 拒绝授权
func AuthorizationDenyHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	if "" == code {
		wirteError(w, -1, errors.New("code is empty"))
		return
	}

	authorizationCode, ok := authorizationEmailMap[code]

	if !ok {
		wirteError(w, -2, errors.New("code is an invalid value"))
		return
	}

	delete(authorizationEmailMap, code)

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		wirteError(w, -3, errors.New("code is an invalid value"))
		return
	}

	fmt.Println("Deny Authorization", code, authorizationCode)

	authorizationInfo.AuthorizationStatus = statusAuthorizationUserDeny

	authorizationMap[authorizationCode] = authorizationInfo

	updateUserTokenTime(authorizationInfo)

	wirteBody(w, 1, "Denied authorization")
}

// AuthorizationUpdateHandler 刷新授权状态
func AuthorizationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	authorizationCode := r.PostFormValue("authCode")
	userToken := r.PostFormValue("userToken")

	if "" == authorizationCode {
		wirteError(w, -1, errors.New("authCode is empty"))
		return
	} else if "" == userToken {
		wirteError(w, -2, errors.New("userToken is empty"))
		return
	}

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		wirteError(w, -3, errors.New("authCode is an invalid value"))
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(authorizationInfo.UserTokenEncode), []byte(userToken))

	if err != nil {
		wirteError(w, -4, errors.New("userToken is an invalid value"))
		return
	}

	timestamp := time.Now().Unix()

	if authorizationInfo.TokenExpireTime < timestamp {
		if authorizationInfo.TokenDeleteTime < timestamp {
			delete(authorizationMap, authorizationCode)
		}
		wirteError(w, -5, errors.New("userToken has expired"))
		return
	}

	authorizationStatus := authorizationInfo.AuthorizationStatus

	if statusAuthorizationUserDeny == authorizationStatus {
		wirteError(w, -6, errors.New("The user denied access"))
		return
	}

	updateUserTokenTime(authorizationInfo)

	wirteBody(w, 1, "update userToken successfully")
}

// AuthorizationStateHandler 获取授权状态
func AuthorizationStateHandler(w http.ResponseWriter, r *http.Request) {
	authorizationCode := r.PostFormValue("authCode")
	userToken := r.PostFormValue("userToken")

	if "" == authorizationCode {
		wirteError(w, -1, errors.New("authCode is empty"))
		return
	} else if "" == userToken {
		wirteError(w, -2, errors.New("userToken is empty"))
		return
	}

	authorizationInfo, ok := authorizationMap[authorizationCode]

	if !ok {
		wirteError(w, -3, errors.New("authCode is an invalid value"))
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(authorizationInfo.UserTokenEncode), []byte(userToken))

	if err != nil {
		wirteError(w, -4, errors.New("userToken is an invalid value"))
		return
	}

	timestamp := time.Now().Unix()

	t1 := strconv.FormatInt(timestamp-authorizationInfo.TokenExpireTime, 10)
	t2 := strconv.FormatInt(timestamp-authorizationInfo.TokenDeleteTime, 10)

	fmt.Println(t1, t2)

	if authorizationInfo.TokenExpireTime < timestamp {
		if authorizationInfo.TokenDeleteTime < timestamp {
			delete(authorizationMap, authorizationCode)
		}
		wirteError(w, -5, errors.New("userToken has expired"))
		return
	}

	authorizationStatus := authorizationInfo.AuthorizationStatus

	wirteBody(w, 1, map[string]int{"authStatus": authorizationStatus})
}

// ServiceRegiesterHandler 注册服务
func ServiceRegiesterHandler(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PostFormValue("serviceId")
	serviceName := r.PostFormValue("serviceName")
	serviceDesc := r.PostFormValue("serviceDesc")

	if "" == serviceID {
		wirteError(w, -1, errors.New("serviceId is empty"))
		return
	} else if "" == serviceName {
		wirteError(w, -2, errors.New("serviceName is empty"))
		return
	}

	var providerRegisterResult, requesterRegisterResult = false, false

	registerResult := make(map[string]string)

	if _, ok := serviceMap[serviceID]; !ok {
		token, encodeToken, err := genServiceToken(serviceID)

		if err != nil {
			fmt.Println("Provider", err)
			wirteError(w, -7, errors.New("This service failed to register"))
			return
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
		wirteBody(w, 1, registerResult)
	} else {
		wirteError(w, -6, errors.New("This service is already registered"))
	}
}

// ServiceUnRegiesterHandler 注销服务
func ServiceUnRegiesterHandler(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PostFormValue("serviceId")
	serviceToken := r.PostFormValue("serviceToken")

	if "" == serviceID {
		wirteError(w, -1, errors.New("serviceId is empty"))
		return
	} else if "" == serviceToken {
		wirteError(w, -2, errors.New("serviceToken is empty"))
		return
	}

	serviceProvider, providerOk := serviceMap[serviceID]

	if providerOk {
		err := bcrypt.CompareHashAndPassword([]byte(serviceProvider.ServiceTokenEncode), []byte(serviceToken))

		if err == nil {
			delete(serviceMap, serviceID)
			wirteBody(w, 1, "This service has been successfully unregistered")
			return
		}
	}

	wirteError(w, -6, errors.New("serviceToken is an invalid value"))
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

	grantURL := "http://" + serverDomain + "/authorization/grant?code=" + emailCode
	denyURL := "http://" + serverDomain + "/authorization/deny?code=" + emailCode

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

func wirteResult(w http.ResponseWriter, code int, data interface{}) error {
	resultJSON, err := json.Marshal(Result{code, data})
	if err != nil {
		return err
	}

	wirteResponse(w, string(resultJSON))

	return nil
}

// 统一错误输出接口
func wirteError(w http.ResponseWriter, code int, err error) {
	fmt.Println(code, ",", err)

	errJSON, err := json.Marshal(Result{code, err.Error()})
	if err != nil {
		wirteResponse(w, "{\"code\": "+strconv.Itoa(code)+",\"data\": \""+err.Error()+"\"}")
	} else {
		wirteResponse(w, string(errJSON))
	}
}

func wirteResponse(w http.ResponseWriter, resp string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Content-Type", "Application/json")             //header的类型
	w.Write([]byte(resp))
}

func wirteBody(w http.ResponseWriter, code int, data interface{}) {
	err := wirteResult(w, code, data)
	if err != nil {
		wirteError(w, -100, err)
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

	SetupMailCredentials("Enter e-mail username: ", "Enter e-mail password: ")

	fmt.Println("Whoam is working")

	authorizationMap = make(map[string]AuthorizationInfo)
	serviceMap = make(map[string]ServiceInfo)
	authorizationEmailMap = make(map[string]string)

	http.HandleFunc("/authorization/request", AuthorizationRequestHandler)
	http.HandleFunc("/authorization/grant", AuthorizationGrantHandler)
	http.HandleFunc("/authorization/deny", AuthorizationDenyHandler)
	http.HandleFunc("/authorization/update", AuthorizationUpdateHandler)
	http.HandleFunc("/authorization/state", AuthorizationStateHandler)
	http.HandleFunc("/serice/register", ServiceRegiesterHandler)
	http.HandleFunc("/serice/unregister", ServiceUnRegiesterHandler)
	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(serverPort), nil)
	fmt.Println(err)
}
