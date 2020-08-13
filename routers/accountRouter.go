package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/DuckBap/Duckbap-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SignUpRouter(router *gin.RouterGroup) {
	router.POST("/sign-up", controllers.SignUp, middlewares.Auth.LoginHandler)
	router.POST("/login", middlewares.Auth.LoginHandler)
}
