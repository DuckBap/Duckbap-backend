package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetSwagger(router *gin.RouterGroup) {
	router.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}