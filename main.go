package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var serverPort int
var serverDomain string

func init() {
	ip, err := ExternalIP()
	if err != nil {
		panic(err)
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
		panic("ServerDomain is empty")
	}

	authorizationMap = make(map[string]AuthorizationInfo)
	serviceMap = make(map[string]ServiceInfo)
	authorizationEmailMap = make(map[string]string)
	userVerificationMap = make(map[string]UserVerification)

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
