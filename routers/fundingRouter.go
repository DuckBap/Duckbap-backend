package routers

import (
	"fmt"
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/DuckBap/Duckbap-backend/middlewares"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SetFundingRouter (router *gin.RouterGroup) {
	router.GET("/main", isLogined)
	//r.GET("/main-login", controllers.ListSelect)
	router.GET("/banner", controllers.BannerSelect)
	router.GET("/:fund-id", controllers.GetFunding)
	router.GET("", controllers.GetFundingList)
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
