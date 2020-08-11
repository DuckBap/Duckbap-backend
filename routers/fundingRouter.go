package routers

import (
	"fmt"
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/DuckBap/Duckbap-backend/middlewares"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SetFundingRouter (r *gin.RouterGroup) {
	r.GET("/main/", isLogined)
	//r.GET("/main-login", controllers.ListSelect)
	r.GET("/banner", controllers.BannerSelect)
}

func isLogined (c *gin.Context) {
	test,err :=middlewares.Auth.GetClaimsFromJWT(c)
	if err == nil {
		stringId := fmt.Sprintf("%v", test["id"])
		id,_ := strconv.Atoi(stringId)
		controllers.ListSelect(c, uint(id))
	} else {
		controllers.NotloginListSelect(c)
	}
}
