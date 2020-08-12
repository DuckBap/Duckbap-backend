package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Entertainment struct {
	ID     uint   `json:"entId"`
	Name   string `json:"name"`
	ImgUrl string `json:"imgUrl"`
}

func EnterList(c *gin.Context) {
	var ents []Entertainment
	configs.DB.Table("entertainments").Find(&ents)
	c.JSON(http.StatusOK, gin.H{
		"data": ents,
	})
}
