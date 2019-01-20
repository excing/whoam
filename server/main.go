package main

import (
	"flag"
	"fmt"
	"strconv"
	"net/http"

  "golang.org/x/net/websocket"
)

const STATUS_USER_WAITING = 0

type UserInfo struct {
	// 用户名
	UserName string
	// 用户 Token
	UserToken string
	// 服务提供方
	ServiceProvider string
	// 服务请求方
	ServiceRequester string
	// 服务授权状态
	AuthorizationStatus int
}

// 服务列表
// 不在列表内的服务，将拒绝请求授权
var ServiceList []string
// 用户 Token 的用户信息，包含授权状态
var UserTokenMap map[string]UserInfo

func AuthorizeRequestAsyncHandler(ws *websocket.Conn) {

}

func AuthorizeRequestHandler(w http.ResponseWriter, r *http.Request) {

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

var ServerPort int

func init() {
	flag.IntVar(&ServerPort, "p", 8030, "file server port")
}

func main() {
	flag.Parse()
	flag.Usage()

	fmt.Println("ServerPort: ", ServerPort)

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