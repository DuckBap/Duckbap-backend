package routers

import (
	"github.com/DuckBap/Duckbap-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetArtistRouter (router *gin.RouterGroup) {
	router.GET("", controllers.ShowArtists)
	router.POST("",controllers.CreateArtist)
}
