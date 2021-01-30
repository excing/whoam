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
	tlpUserLogin = "login.html"
	tlpUserOAuth = "oauth.html"

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
	box := packr.NewBox("./html")

	for _, v := range box.List() {
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

	authorized := router.Group("/")
	authorized.Use(AuthRequired)
	{
		// Web page
		authorized.GET("/user/login", handle(loginEndpoint))
		authorized.GET("/user/oauth", handle(oauthEndpoint))
	}

	v1 := router.Group("/api/v1")
	{
		mainRouter := v1.Group("/user/main")
		{
			mainRouter.POST("/code", handle(PostMainCode))
			mainRouter.POST("/auth", handle(PostMainAuth))
		}

		oauthRouter := v1.Group("/user/oauth")
		{
			oauthRouter.POST("/auth", handle(PostUserOAuthAuth))
			oauthRouter.POST("/refresh", handle(PostUserOAuthRefresh))

			oauthRouter.GET("/base", handle(GetUser))
			oauthRouter.GET("/token", handle(GetOAuthCode))
			oauthRouter.GET("/state", handle(GetOAuthState))
		}

		serviceRouter := v1.Group("/service")
		{
			serviceRouter.POST("/", handle(PostService))
		}
	}

	router.Run(":" + strconv.Itoa(config.Port))
}
