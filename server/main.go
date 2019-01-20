package main

import (
	"flag"
	"fmt"
	"strconv"
	"net/http"

  "golang.org/x/net/websocket"
)

// 授权状态：等待中
const STATUS_AUTHORIZATION_WAITING = 0
// 授权状态：同意授权
const STATUS_AUTHORIZATION_USER_GRANT = 1
// 授权状态：拒绝授权
const STATUS_AUTHORIZATION_USER_Deny = -1

type ServiceInfo struct {
	// 服务 ID
	ServiceId string
	// 服务名
	ServiceName string
	// 服务详细，可选
	ServiceDesc string
}

type UserInfo struct {
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

// 服务列表
// 不在列表内的服务，将拒绝请求授权
var ServiceList []string
// 用户授权列表。key 为 UserToken，value 为用户 Token 的用户信息，包含授权状态
var UserTokenMap map[string]UserInfo
// 已注册服务列表。key 为 ServiceId，value 为 ServiceInfo
var ServiceMap map[string]ServiceInfo

func AuthorizeRequestAsyncHandler(ws *websocket.Conn) {

}

// /authorize/request?username=&providerId=&requesterId=
func AuthorizeRequestHandler(w http.ResponseWriter, r *http.Request) {
	UserName := r.PostFormValue("username")
	ServiceProviderId := r.PostFormValue("providerId")
	ServiceRequesterId := r.PostFormValue("requesterId")

	UserToken, UserTokenHash := hash(UserName, ServiceProviderId, ServiceRequesterId)

	userInfo := UserInfo{
		UserName,
		UserTokenHash,
		ServiceProviderId,
		ServiceRequesterId,
		STATUS_AUTHORIZATION_WAITING,
	}

	UserTokenMap[UserToken] = userInfo

	fmt.Println(UserTokenMap)

	go sendAuthorizationEmail(UserName, UserToken, ServiceProviderId, ServiceRequesterId)

	w.Write([]byte(UserToken))
}

func AuthorizeGrantHandler(w http.ResponseWriter, r *http.Request) {

}

func AuthorizeDenyHandler(w http.ResponseWriter, r *http.Request) {

}

func AuthorizeStateHandler(w http.ResponseWriter, r *http.Request) {

}

func ServiceRegiesterHandler(w http.ResponseWriter, r *http.Request) {

}

func ServiceUnRegiesterHandler(w http.ResponseWriter, r *http.Request) {

}

func hash(UserName string, ServiceProviderId string, ServiceRequesterId string) (string, string) {
	return "ha", "hahahahaha"
}

func sendAuthorizationEmail(UserName string, UserToken string, ServiceProviderId string, ServiceRequesterId string) {

}

var ServerPort int

func init() {
	flag.IntVar(&ServerPort, "p", 8030, "file server port")
}

func main() {
	flag.Parse()
	flag.Usage()

	fmt.Println("ServerPort: ", ServerPort)

	UserTokenMap = make(map[string]UserInfo)
	ServiceMap = make(map[string]ServiceInfo)

	http.Handle("/ws/authorize/request", websocket.Handler(AuthorizeRequestAsyncHandler))
	http.HandleFunc("/authorize/request", AuthorizeRequestHandler)
	http.HandleFunc("/authorize/grant", AuthorizeGrantHandler)
	http.HandleFunc("/authorize/deny", AuthorizeDenyHandler)
	http.HandleFunc("/authorize/state", AuthorizeStateHandler)
	http.HandleFunc("/serice/register", ServiceRegiesterHandler)
	http.HandleFunc("/serice/unregister", ServiceUnRegiesterHandler)
	err := http.ListenAndServe("0.0.0.0:" + strconv.Itoa(ServerPort), nil)
	fmt.Println(err)
}