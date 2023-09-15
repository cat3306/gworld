package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func StartHttpServer(port string) {
	ginEngine := gin.New()
	//engine.TrustedPlatform = "X-Client-IP"
	ginEngine.Use(gin.Recovery())
	apiGroup := ginEngine.Group("/api/game")
	{
		gateApi := apiGroup.Group("/gate/")
		gateApi.POST("info", GateInfo)
	}
	err := ginEngine.Run(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		panic(err)
	}
}
