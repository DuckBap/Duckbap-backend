package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetEnterRouter(router *gin.RouterGroup) {
	router.GET("", controllers.EnterList)
	router.POST("", controllers.CreateEnter)
}
