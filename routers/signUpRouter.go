package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SignUpRouter(router *gin.RouterGroup) {
	router.GET("/sign-up", controllers.ShowArtists)
	router.POST("/sign-up", controllers.SignUp)
}
