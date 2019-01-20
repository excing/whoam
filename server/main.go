package main

import (
	"flag"
	"fmt"
	"strconv"
	"errors"
	"net/http"
	"encoding/json"

  "golang.org/x/net/websocket"
)

// 授权状态：等待中
const STATUS_AUTHORIZATION_WAITING = 0
// 授权状态：同意授权
const STATUS_AUTHORIZATION_USER_GRANT = 1
// 授权状态：拒绝授权
const STATUS_AUTHORIZATION_USER_Deny = -1

// 服务类型：内容提供者
const TYPE_SERVICE_PROVIDER = 1
// 服务类型：内容消费者
const TYPE_SERVICE_REQUESTER = 2
// 服务类型最小值
const TYPE_SERVICE_MIN_VALUE = 1
// 服务类型最大值
const TYPE_SERVICE_MAX_VALUE = 3

type ServiceInfo struct {
	// 服务 ID
	ServiceId string
	// 服务名
	ServiceName string
	// 服务详细，可选
	ServiceDesc string
}

// 认证请求信息
// 请求用户将 ServiceProvider 服务授权给 ServiceRequester
type AuthorizationInfo struct {
	// 用户名
	UserName string
	// 用户 Token SHA256 值
	UserTokenHash string
	// 提供方服务ID
	ServiceProviderId string
	// 请求方服务ID
	ServiceRequesterId string
	// 服务授权状态
	// see: STATUS_AUTHORIZATION_WAITING, STATUS_AUTHORIZATION_USER_GRANT, STATUS_AUTHORIZATION_USER_Deny
	AuthorizationStatus int
}

// 返回消息结构
type Result struct {
	Code int							`json:"code"`
	Data interface{}			`json:"data"`
}

// 服务列表
// 不在列表内的服务, 将拒绝请求授权
var serviceList []string
// 授权请求列表. key 为 RequestId, value 为用户 Token 的用户信息, 包含授权状态
var authorizationMap map[string]AuthorizationInfo
// 服务提供方列表. key 为 ServiceId, value 为 ServiceInfo
var serviceProviderMap map[string]ServiceInfo
// 服务请求方列表. key 为 ServiceId, value 为 ServiceInfo
var serviceRequesterMap map[string]ServiceInfo

func AuthorizationRequestAsyncHandler(ws *websocket.Conn) {

}

// /authorization/request?username=&providerId=&requesterId=
func AuthorizationRequestHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.PostFormValue("username")
	serviceProviderId := r.PostFormValue("providerId")
	serviceRequesterId := r.PostFormValue("requesterId")

	if "" == userName {
		wirteError(w, -1, errors.New("username is empty"))
		return
	} else if "" == serviceProviderId {
		wirteError(w, -2, errors.New("providerId is empty"))
		return
	} else if "" == serviceRequesterId {
		wirteError(w, -3, errors.New("requesterId is empty"))
		return
	} else if serviceProviderId == serviceRequesterId {
		wirteError(w, -6, errors.New("providerId equals requesterId"))
		return
	}

	serviceProvider, serviceProviderIdOk := serviceProviderMap[serviceProviderId]
	serviceRequester, serviceRequesterIdOk := serviceRequesterMap[serviceRequesterId]

	if !serviceProviderIdOk {
		wirteError(w, -4, errors.New("Can not find the service pointed to by providerId"))
		return
	} else if !serviceRequesterIdOk {
		wirteError(w, -5, errors.New("Can not find the service pointed to by requesterId"))
		return
	}

	authorizationCode, userToken, userTokenHash := genCodeAndToken(userName, serviceProvider, serviceRequester)

	authorizationMap[authorizationCode] = AuthorizationInfo {
		userName,
		userTokenHash,
		serviceProviderId,
		serviceRequesterId,
		STATUS_AUTHORIZATION_WAITING,
	}

	go sendAuthorizationEmail(authorizationCode, authorizationMap[authorizationCode], userToken)

	data := make(map[string]string)
	data["authCode"] = authorizationCode
	data["userToken"] = userToken

	wirteBody(w, 1, data)
}

func AuthorizationGrantHandler(w http.ResponseWriter, r *http.Request) {

}

func AuthorizationDenyHandler(w http.ResponseWriter, r *http.Request) {

}

func AuthorizationStateHandler(w http.ResponseWriter, r *http.Request) {

}

func ServiceRegiesterHandler(w http.ResponseWriter, r *http.Request) {
	serviceId := r.PostFormValue("serviceId")
	serviceName := r.PostFormValue("serviceName")
	serviceDesc := r.PostFormValue("serviceDesc")
	serviceTypeStr := r.PostFormValue("serviceType")

	if "" == serviceId {
		wirteError(w, -1, errors.New("serviceId is empty"))
		return
	} else if "" == serviceName {
		wirteError(w, -2, errors.New("serviceName is empty"))
		return
	} else if "" == serviceTypeStr {
		wirteError(w, -3, errors.New("serviceType is empty"))
		return
	}

	serviceType, err := strconv.Atoi(serviceTypeStr)

	if err != nil {	
		wirteError(w, -4, errors.New("serviceType is an invalid value"))
		return
	} else if serviceType < TYPE_SERVICE_MIN_VALUE || TYPE_SERVICE_MAX_VALUE < serviceType {
		wirteError(w, -5, errors.New("serviceType is a wrong value"))
		return
	}

	var providerRegisterResult, requesterRegisterResult = false, false

	if TYPE_SERVICE_PROVIDER == (serviceType & TYPE_SERVICE_PROVIDER) {
		if _, ok := serviceProviderMap[serviceId]; !ok {
			serviceProviderMap[serviceId] = ServiceInfo {
				serviceId,
				serviceName,
				serviceDesc,
			}
			providerRegisterResult = true
		}
	}

	if TYPE_SERVICE_REQUESTER == (serviceType & TYPE_SERVICE_REQUESTER) {
		if _, ok := serviceRequesterMap[serviceId]; !ok {
			serviceRequesterMap[serviceId] = ServiceInfo {
				serviceId,
				serviceName,
				serviceDesc,
			}
			requesterRegisterResult = true
		}
	}

	if providerRegisterResult || requesterRegisterResult {
		wirteBody(w, 1, "This service is registered successfully")
	} else {
		wirteError(w, -6, errors.New("This service is already registered"))
	}
}

func ServiceUnRegiesterHandler(w http.ResponseWriter, r *http.Request) {
	serviceId := r.PostFormValue("serviceId")

	if "" == serviceId {
		wirteError(w, -1, errors.New("serviceId is empty"))
		return
	}

	delete (serviceProviderMap, serviceId)
	delete (serviceRequesterMap, serviceId)

	wirteBody(w, 1, "This service has been successfully unregistered")
}

func genCodeAndToken(userName string, serviceProvider ServiceInfo, serviceRequester ServiceInfo) (string, string, string) {
	return "a123", "ha", "hahahahaha"
}

func sendAuthorizationEmail(authorizationCode string, authorizationInfo AuthorizationInfo, userToken string) {

}

func wirteResult(w http.ResponseWriter, code int, data interface{}) error {
	resultJson, err := json.Marshal( Result{code, data} )
	if err != nil {
		return err
	}

	wirteResponse(w, string(resultJson))

	return nil
}

// 统一错误输出接口
func wirteError(w http.ResponseWriter, code int, err error) {
	fmt.Println(code , ",", err)

	errJson, err := json.Marshal( Result{code, err.Error()} )
	if err != nil {
		wirteResponse(w, "{\"code\": " + strconv.Itoa(code) + ",\"data\": \"" + err.Error() + "\"}")
	} else {
		wirteResponse(w, string(errJson))
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

var ServerPort int

func init() {
	flag.IntVar(&ServerPort, "p", 8030, "file server port")
}

func main() {
	flag.Parse()
	flag.Usage()

	fmt.Println("ServerPort: ", ServerPort)

	authorizationMap = make(map[string]AuthorizationInfo)
	serviceProviderMap = make(map[string]ServiceInfo)
	serviceRequesterMap = make(map[string]ServiceInfo)

	http.Handle("/ws/authorization/request", websocket.Handler(AuthorizationRequestAsyncHandler))
	http.HandleFunc("/authorization/request", AuthorizationRequestHandler)
	http.HandleFunc("/authorization/grant", AuthorizationGrantHandler)
	http.HandleFunc("/authorization/deny", AuthorizationDenyHandler)
	http.HandleFunc("/authorization/state", AuthorizationStateHandler)
	http.HandleFunc("/serice/register", ServiceRegiesterHandler)
	http.HandleFunc("/serice/unregister", ServiceUnRegiesterHandler)
	err := http.ListenAndServe("0.0.0.0:" + strconv.Itoa(ServerPort), nil)
	fmt.Println(err)
}