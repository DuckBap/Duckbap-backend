package routers

import (
	"fmt"
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetFundingRouter (r *gin.RouterGroup) {
	r.GET("/main/", isLogined)
	//r.GET("/main-login", controllers.ListSelect)
	r.GET("/banner", controllers.BannerSelect)
}

func isLogined (c *gin.Context) {
	_, ok := c.Get("user")
	if ok {
		fmt.Println("ok")
		controllers.ListSelect(c)
	} else {
		fmt.Println("failed")
		controllers.NotloginListSelect(c)
	}
}