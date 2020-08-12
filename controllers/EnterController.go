package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func	EnterList(c *gin.Context) {
	var ents []models.Entertainment
	configs.DB.Table("entertainments").Find(&ents)
	c.JSON(http.StatusOK,gin.H{
		"data": ents,
	})
}