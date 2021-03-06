package controllers

import (
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

// @Summary 아티스트 리스트
// @Description ## <br>아티스트 리스트를 반환합니다.
// @Description <br>
// @Description # 쿼리스트링이 존재하지 않는 경우
// @Description 1. 모든 아티스트를 반환합니다.<br>
// @Description # <br>쿼리스트링이 존재하는 경우
// @Description 1. 쿼리스트링을 조건으로 필터링 된 아티스트를 반환합니다.<br>
// @Description 1. 회사에 속한 아티스트들 /v1/artists?ent-id=()
// @Accept  json
// @Produce  json
// @Router /artists [get]
// @Success 200 {array} OutputArtistList
func ShowArtists(c *gin.Context) {
	id, _ := c.GetQuery("ent-id")
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

func CreateArtist(c *gin.Context) {
	var artist models.Artist
	c.BindJSON(&artist)

	configs.DB.Create(&artist)
	c.JSON(http.StatusOK, artist)
}
