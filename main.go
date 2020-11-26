package main

import (
	"html/template"
	"strconv"
	"time"

	"github.com/excing/goflag"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
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

const (
	tlpUserLogin = "userOAuthLogin.html"
	tlpUserOAuth = "userOAuth.html"

	tlpFaviconSVG = "favicon.svg"

	// MainServiceID main servvice id
	MainServiceID = "whoam.xyz"
)

var config Config
var db *gorm.DB

func init() {
	port := 8030
	ip, err := ExternalIP()
	if err != nil {
		panic(err)
	}

	config = Config{port, ip + ":" + strconv.Itoa(port), "test.db", false}

	goflag.Var(&config)
}

func main() {
	goflag.Parse("config", "Configuration file path")

	time.FixedZone("CST", 8*3600)

	var err error
	db, err = gorm.Open(sqlite.Open(config.Db), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	initUser()

	tmpl := template.New("user")
	box := packr.NewBox("./web")
	htmls := []string{
		tlpUserLogin,
		tlpUserOAuth,
		tlpFaviconSVG,
	}

	for _, v := range htmls {
		indexTmpl := tmpl.New(v)
		data, _ := box.FindString(v)
		indexTmpl.Parse(data)
	}

	router := gin.Default()
	router.SetHTMLTemplate(tmpl)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.StaticFS("/favicon_io", packr.NewBox("./favicon_io"))

	router.GET("/user/login", inout(PageUserLogin))
	router.GET("/user/oauth", inout(PageUserOAuth))

	apiV1 := router.Group("/api/v1")

	apiV1.POST("/user/main/code", inout(PostMainCode))
	apiV1.POST("/user/main/auth", inout(PostUserAuth))

	apiV1.POST("/user/oauth/auth", inout(PostUserOAuthAuth))
	apiV1.GET("/user/oauth/token", inout(GetOAuthCode))
	apiV1.GET("/user/oauth/state", inout(GetOAuthState))

	apiV1.POST("/servicer", inout(PostServicer))
	apiV1.DELETE("/servicer", inout(DeleteServicer))
	apiV1.POST("/service/method", inout(PostServiceMethod))

	router.Run(":" + strconv.Itoa(config.Port))
}
