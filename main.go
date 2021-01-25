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

	tlpFaviconSVG = "favicon.svg"

	// MainServiceID main servvice id
	MainServiceID = "whoam.xyz"
)

var config Config
var ctx context.Context
var client *ent.Client

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

	InitRAS()
	InitUser()
	InitService()

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

	v1 := router.Group("/api/v1")

	v1.POST("/user/main/code", inout(PostMainCode))
	v1.POST("/user/main/auth", inout(PostUserAuth))

	v1.POST("/user/oauth/auth", inout(PostUserOAuthAuth))
	v1.GET("/user/oauth/token", inout(GetOAuthCode))
	v1.GET("/user/oauth/state", inout(GetOAuthState))
	v1.POST("/user/oauth/refresh", inout(PostUserOAuthRefresh))

	v1.POST("/service", inout(PostService))
	v1.POST("/service/:id/method", inout(PostServiceMethod))
	v1.POST("/service/:id/permission", inout(PostServicePermission))

	v1.POST("/article/new", inout(NewArticle))
	v1.GET("/article/:id", inout(GetArticle))
	v1.GET("/articles", inout(GetArticles))
	v1.POST("/accord/new", inout(NewAccord))
	v1.GET("/accord/:id", inout(GetAccord))
	v1.GET("/accords", inout(GetAccords))
	v1.GET("/accord/:id/articles", inout(GetAccordArticles))

	v1.POST("/ras/new", inout(NewRAS))
	v1.GET("/ras/user/:userId", inout(GetRAS))
	v1.POST("/ras/vote", inout(VoteRAS))
	v1.GET("/votes/rasId/:rasId", inout(GetRasVotes))
	v1.GET("/votes/post/:postUri", inout(GetPostVotes))

	router.Run(":" + strconv.Itoa(config.Port))
}
