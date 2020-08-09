package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetFundingRouter (r *gin.RouterGroup) {
	r.GET("/main_login", controllers.ListSelect)
	r.GET("/banner", controllers.BannerSelect)
}