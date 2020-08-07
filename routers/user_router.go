package routers

import (
	"github.com/DuckBap/Duckbap-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUserRouters(r *gin.RouterGroup) {
	r.POST("/login", middlewares.Auth.LoginHandler)
	//r.GET("/me", middlewares.Auth.MiddlewareFunc(), <HANDLER>)
}
