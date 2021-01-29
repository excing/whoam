package main

import (
	"context"
	"html/template"
	"strconv"
	"time"

	"github.com/excing/goflag"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	_ "github.com/mattn/go-sqlite3"
	"whoam.xyz/ent"
)

// Config 配置文件信息
type Config struct {
	Port  int    `flag:"Authorization server port"`
	Db    string `flag:"Authorization database file path"`
	Debug bool   `flag:"Is Debug mode"`
}

const (
	tlpUserLogin = "userLogin.html"
	tlpUserOAuth = "userOAuth.html"

	// MainServiceID main servvice id
	MainServiceID = "whoam.xyz"
)

var config Config
var ctx context.Context
var client *ent.Client
var router *gin.Engine

func init() {
	config = Config{Port: 8030, Db: "test.db", Debug: false}

	goflag.Var(&config)
}

func main() {
	goflag.Parse("config", "Configuration file path")

	time.FixedZone("CST", 8*3600)

	var err error

	opts := []ent.Option{}
	if config.Debug {
		opts = append(opts, ent.Debug())
	}
	if config.Db == "" {
		client, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
	} else {
		client, err = ent.Open("sqlite3", "file:"+config.Db+"?_fk=1", opts...)
	}

	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	defer client.Close()

	ctx = context.Background()
	if err = client.Schema.Create(ctx); err != nil {
		panic("failed to create schema: " + err.Error())
	}

	InitUser()
	InitService()

	tmpl := template.New("user")
	box := packr.NewBox("./web")
	htmls := []string{
		tlpUserLogin,
		tlpUserOAuth,
	}

	for _, v := range htmls {
		indexTmpl := tmpl.New(v)
		data, _ := box.FindString(v)
		indexTmpl.Parse(data)
	}

	router = gin.Default()
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

	router.GET("/user/login", authorizeUser, handle(PageUserLogin))
	router.GET("/user/oauth", authorizeUser, handle(PageUserOAuth))

	v1 := router.Group("/api/v1")

	v1.POST("/user/main/code", handle(PostMainCode))
	v1.POST("/user/main/auth", handle(PostUserAuth))

	v1.POST("/user/oauth/auth", handle(PostUserOAuthAuth))
	v1.POST("/user/oauth/refresh", handle(PostUserOAuthRefresh))

	v1.GET("/user", handle(GetUser))
	v1.GET("/user/oauth/token", handle(GetOAuthCode))
	v1.GET("/user/oauth/state", handle(GetOAuthState))

	v1.POST("/service", handle(PostService))
	v1.POST("/service/method", handle(PostServiceMethod))
	v1.POST("/service/permission", handle(PostServicePermission))

	router.Run(":" + strconv.Itoa(config.Port))
}
