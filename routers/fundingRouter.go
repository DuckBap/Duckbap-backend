package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetFundingUrls(r *gin.RouterGroup){
	r.GET("", controllers.MainFundingList)
	r.GET("/banner", controllers.BannerFundingList)
}
