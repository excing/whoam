package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Config 配置文件信息
type Config struct {
	Port   int    `flag:"Authorization server port"`
	Domain string `flag:"Authorization server domain"`
	Db     string `flag:"Authorization database file path"`
	Debug  bool   `flag:"Is Debug mode"`
}

var config Config
var db *gorm.DB

func init() {
	port := 8030
	ip, err := ExternalIP()
	if err != nil {
		panic(err)
	}

	config = Config{port, ip + ":" + strconv.Itoa(port), "test.db", false}

	FlagVar(&config)
}

func main() {
	FlagParse("config", "Configuration file path")

	time.FixedZone("CST", 8*3600)

	var err error
	db, err = gorm.Open(sqlite.Open(config.Db), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	initUser()

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

	router.Run(":" + strconv.Itoa(config.Port))
}
