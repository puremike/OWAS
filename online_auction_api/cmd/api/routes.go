package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func (app *application) routes() http.Handler {

	g := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	api := g.Group("/api/v1")
	{
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		api.GET("/health", app.health)
	}

	return g
}
