package controllers

import (
	"fmt"
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OutputArtistList struct {
	ID              uint   `json:"artistId"`
	Name            string `json:"artistName"`
	ImgUrl          string `json:"artistImgUrl"`
	EntertainmentID uint   `json:"entId"`
}

func ShowArtists(c *gin.Context) {
	id, _ := c.GetQuery("ent-id")
	fmt.Println(id)
	entId, _ := strconv.Atoi(id)
	if entId != 0 {
		EnterArtistList(c, entId)
	} else {
		ArtistList(c)
	}
}

func EnterArtistList(c *gin.Context, id int) {
	var list []OutputArtistList
	configs.DB.Table("artists").Where("entertainment_id = ?", id).Find(&list)
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func ArtistList(c *gin.Context) {
	var list []OutputArtistList
	configs.DB.Table("artists").Find(&list)
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func	CreateArtist(c *gin.Context) {
	var artist models.Artist
	c.BindJSON(&artist)

	configs.DB.Create(&artist)
	c.JSON(http.StatusOK, artist)
}