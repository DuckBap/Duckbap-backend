package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Entertainment struct {
	ID     uint   `json:"entId"`
	Name   string `json:"name"`
	ImgUrl string `json:"imgUrl"`
}

// @Summary 엔터테인먼트 리스트
// @Description <br>엔터테인먼트 리스트를 반환합니다.<br>
// @Accept  json
// @Produce  json
// @Router /ents [get]
// @Success 200 {array} Entertainment
func EnterList(c *gin.Context) {
	var ents []Entertainment
	configs.DB.Table("entertainments").Find(&ents)
	c.JSON(http.StatusOK, gin.H{
		"data": ents,
	})
}

func CreateEnter(c *gin.Context) {
	var ent models.Entertainment
	c.BindJSON(&ent)

	configs.DB.Create(&ent)
	c.JSON(http.StatusOK, ent)
}
