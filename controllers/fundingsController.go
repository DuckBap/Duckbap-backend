package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/gin-gonic/gin"
	"math"
)

type mainFunding struct {
	ID           uint
	Name         string
	MainImgUrl   string
	AttainmentRate float64
	TargetAmount uint `json:"-"`
	SalesAmount  uint `json:"-"`
	ArtistID     uint `json:"-"`
}

type bannerFunding struct {
	ID uint
	Name string
	MainImgUrl string
}

func BannerFundingList (c *gin.Context) {
	var fundings []bannerFunding

	configs.DB.Table("fundings").Order("sales_amount desc").Limit(5).Find(&fundings)
	c.JSON(200, fundings)
}

func MainFundingList(c *gin.Context) {
	var fundings []mainFunding
	var res []mainFunding

	userID := 1

	artistsInBookmarks := configs.DB.Table("bookmarks").Select("artist_id").Where("user_id = ?", userID)
	favoriteArtist := configs.DB.Table("users").Select("favorite_artist").Where("id = ?", userID)
	configs.DB.Table("fundings").Where("artist_id in (?) || artist_id in (?)", artistsInBookmarks, favoriteArtist).Order("sales_amount desc").Find(&fundings)
	artistIDs := getArtists(fundings)

	for _, artistID := range artistIDs {
		limit := math.Ceil(8 / float64(len(artistIDs)))
		for _, funding := range fundings {
			if len(res) >= 8 {
				break
			}
			if funding.ArtistID == artistID && !isFundingDuplicated(funding, res) && limit > 0 {
				res = append(res, funding)
				limit--
			}
		}
	}
	if len(res) < 8 {
		configs.DB.Table("fundings").Order("sales_amount desc").Limit(8).Find(&fundings)
		for _, funding := range fundings {
			if !isFundingDuplicated(funding, res) {
				res = append(res, funding)
			}
			if len(res) >= 8 {
				break
			}
		}
	}
	setAttainment(getAddressList(&res))
	c.JSON(200, res)
}

func getAddressList(slice *[]mainFunding) (addressList []*mainFunding) {

	for i, _ := range *slice {
		addressList = append(addressList, &((*slice)[i]))
	}
	return
}

func setAttainment(res []*mainFunding) {
	for _, funding := range res {
		funding.AttainmentRate = math.Round(float64(funding.SalesAmount) / float64(funding.TargetAmount) * 100)
	}
}

func isFundingDuplicated(target mainFunding, values []mainFunding) bool {
	for _, v := range values {
		if v.ID == target.ID {
			return true
		}
	}
	return false
}

func isNumDuplicated(target uint, slice []uint) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func getArtists(fundings []mainFunding) (artistIDs []uint) {

	for _, funding := range fundings {
		if !isNumDuplicated(funding.ArtistID, artistIDs) {
			artistIDs = append(artistIDs, funding.ArtistID)
		}
	}
	return
}

//중복된 펀딩이 아닐 것
//연예인 limit을 초과하지 않을 것 --- ok
//8개를 넘어가지 않을 것
