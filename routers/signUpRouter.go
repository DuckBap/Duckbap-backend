package routers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
)

func SignUpRouter(router *gin.RouterGroup) {
	router.GET("/mi", func (c *gin.Context) {
		configs.DB.AutoMigrate(&models.User{}, &models.Funding{}, &models.FundingImg{},&models.Artist{}, &models.Receipt{},  &models.Entertainment{})
	})
	router.GET("/sign-up", controllers.ShowArtists)
	router.POST("/sign-up", controllers.SignUp)
}
