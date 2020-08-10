package routers

import (
	controller "github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetFundingRouter(router *gin.RouterGroup) {
	// group name : /fundings
	//router.POST("/", controller.CreateFunding)
	router.GET("/:fund_id", controller.GetFunding)
	router.GET("/", controller.GetFundingList)
}