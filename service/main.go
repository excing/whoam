package main

import (
	"github.com/gin-gonic/gin"
	"whoam.xyz/api"
)

func main() {
	router := gin.Default()
	v1 := router.Group("/v1")

	v1.GET("/ws/authorization/request", api.Output(api.WSAuthorizationRequest))
	// v1.POST("/authorization/request", postAuthorizationRequest)
	// v1.POST("/authorization/grant", postAuthorizationGrant)
	// v1.POST("/authorization/deny", postAuthorizationDeny)
	// v1.POST("/authorization/update", postAuthorizationUpdate)
	// v1.POST("/authorization/state", postAuthorizationState)
	// v1.POST("/serice/register", postSericeRegister)
	// v1.POST("/serice/unregister", postSericeUnregister)

	router.Run(":12301")
}
