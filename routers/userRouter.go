package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/DuckBap/Duckbap-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUserRouters(router *gin.RouterGroup) {
	router.GET("/me", middlewares.Auth.MiddlewareFunc(), controllers.GetMe)
	router.POST("/bookmark", controllers.CreateBookmark)
	router.POST("/receipt", controllers.CreateReceipts)
}
