package routers

import (
	"github.com/DuckBap/Duckbap-backend/middlewares"
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"

)

func SetUserRouters(router *gin.RouterGroup) {
	router.GET("/me", middlewares.Auth.MiddlewareFunc(), controllers.GetMe)
	//r.GET("/me", middlewares.Auth.MiddlewareFunc(), <HANDLER>)
}
