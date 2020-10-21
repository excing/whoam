package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIMessage API 消息体
type APIMessage struct {
	Code int         `json:"code"`
	Type string      `json:"-"`
	Data interface{} `json:"data"`
}

func (msg *APIMessage) Error() string {
	return fmt.Sprintf("%v", msg.Data)
}

func stringAPI(code int, data string) APIMessage {
	return APIMessage{Code: code, Type: "string", Data: data}
}

func jsonAPI(code int, data interface{}) APIMessage {
	return APIMessage{Code: code, Type: "json", Data: data}
}

func binaryAPI(code int, data string) APIMessage {
	return APIMessage{Code: code, Type: "binary", Data: data}
}

// Output 统一输出
func Output(fn func(*gin.Context) APIMessage) gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := fn(c)

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		switch msg.Type {
		case "string":
			c.String(msg.Code, msg.Data.(string))
		case "json":
			c.JSON(msg.Code, msg.Data)
		case "binary":
			c.File(msg.Data.(string))
		}
	}
}

// WSAuthorizationRequest 认证请求
func WSAuthorizationRequest(c *gin.Context) APIMessage {
	return jsonAPI(http.StatusOK, 0)
}
