package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
)

type listFunding struct {
	ID uint
	SellerID uint
	Name string
	Price uint
	TargetAmount uint
	MainImgUrl string
	ArtistID uint
	SalesAmount uint
}

type bannerFunding struct {
	ID uint
	Name string
	MainImgUrl string
	ArtistID uint
}

type bookmarks struct {
	ArtistID uint
}

func bannerSelect (c *gin.Context) {
	var fundings []bannerFunding

	configs.DB.Table("fundings").Order("sales_amount desc").Limit(5).Find(&fundings)
	c.JSON(200, fundings)
}

func listSelect (c *gin.Context) {
	var fundings []listFunding
	var tmp []listFunding

	bookmark := findBookmarks(c)
	limit := int(math.Ceil(8.0/float64(len(bookmark))))
	for i:=0;i<len(bookmark);i++ {
		configs.DB.Table("fundings").Where("artist_id = ?", bookmark[i]).Order("sales_amount desc").Limit(limit).Find(&tmp)
		fundings = append(fundings, tmp...)
	}
	sort.Slice(fundings, func (i, j int) bool {
		return fundings[i].SalesAmount > fundings[j].SalesAmount
	})
	if len(fundings) < 8 {
		configs.DB.Table("fundings").Not("artist_id", bookmark).Order("sales_amount desc").Limit(8 - len(fundings)).Find(&tmp)
		fundings = append(fundings, tmp...)
	}
	c.JSON(200, fundings)
}

func findBookmarks (c *gin.Context) []bookmarks {
	id := c.Params.ByName("id")
	var bookmark []bookmarks
	var favorite bookmarks

	configs.DB.Table("bookmarks").Where("user_id = ?", id).Order("artist_id").Find(&bookmark)
	configs.DB.Table("users").Select("favorite_artist").Where("user_id = ?", id).Find(&favorite)
	bookmark = append(bookmark, favorite)
	return bookmark
}