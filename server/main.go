package main

import (
	_ "event_trigger/docs"
	"event_trigger/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	apiVersion = "v1"
)

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	a := routes.New(apiVersion)

	a.Start(a.Engine)
}
