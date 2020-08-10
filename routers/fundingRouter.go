package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetFundingRouter (r *gin.RouterGroup) {
	r.GET("/main/", isLogined)
	r.GET("/banner", controllers.BannerSelect)
}

func isLogined (c *gin.Context) {
	_, ok := c.Get("user")
	if ok {
		controllers.ListSelect(c)
	} else {
		controllers.NotloginListSelect(c)
	}
}