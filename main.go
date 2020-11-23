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
	tlpUserOAuthLogin = "userOAuthLogin.html"
	tlpFaviconSVG     = "favicon.svg"

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
		tlpUserOAuthLogin,
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
	router.GET("/favicon.ico", inout(func(c *Context) error {
		return c.OkHTML(tlpFaviconSVG, nil)
	}))

	v1 := router.Group("/v1")
	v1.POST("/user/main/code", inout(PostMainCode))
	v1.POST("/user/main/auth", inout(PostUserAuth))
	v1.GET("/user/main/state", inout(GetUserState))

	v1.GET("/user/oauth/login", inout(PostUserOAuthLogin))
	v1.POST("/user/oauth/auth", inout(PostUserOAuthAuth))
	v1.GET("/user/oauth/token", inout(GetOAuthCode))
	v1.POST("/user/oauth/state", inout(GetOAuthState))

	v1.POST("/servicer", inout(PostServicer))
	v1.DELETE("/servicer", inout(DeleteServicer))

	router.Run(":" + strconv.Itoa(config.Port))
}
