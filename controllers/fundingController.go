package controllers

import (
	"database/sql"
	"errors"
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

//var db = configs.DB
/*
## 제품 사진

- 5장 - fund_img

## 제품 정보

- 이름 - fund/
- 한 세트 당 가격 - fund/
- 펀딩 목표량 - fund/
- 현재 달성량 (주문량) - fund/
- 달성률 - fund/
- 시작일 - fund/
- 마감일 - fund/
- 디데이 - fund/
- 판매자 - fund/
- 아티스트 이름 - artist/

## 상세 정보

- 이미지
*/

type fundingResBody struct {
	ID              uint      `json:"id"`
	NickName        string    `json:"sellerName"`
	Name            string    `json:"fundName"`
	Price           uint      `json:"price"`
	TargetAmount    uint      `json:"targetAmount"`
	SalesAmount     uint      `json:"salesAmount"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	ArtistName      string    `json:"artistName"`
	AchievementRate float64   `json:"achievementRate"` //salesAmount / Price
	Dday            uint      `json:"dDay"`
	FundingImgUrls  []string  `json:"fundingImgUrls"`
	DetailedImgUrl  string    `json:"detailedImgUrl"`
}

//test
//func CreateFunding(c *gin.Context) {
//	fund := models.Funding{
//		SellerID: 1,
//		Name: "아이 굿즈3",
//		Price: 4000,
//		TargetAmount: 50000,
//		MainImgUrl: "이미지4",
//		ArtistID: 2,
//		StartDate: time.Now(),
//		EndDate: time.Now().Add(24 * time.Hour),
//	}
//	configs.DB.Create(&fund)
//	c.JSON(http.StatusOK, gin.H{
//		"msg": "create funding",
//		"funding": fund,
//	})
//}

func GetFunding(c *gin.Context) {
	fundID := c.Param("fund-id")

	body, err := setFundingBody(fundID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "no funding",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":     "get funding",
			"funding": body,
		})
	}

}

func setFundingBody(fundID string) (*fundingResBody, error) {
	body := fundingResBody{}
	var titleImg string

	sqlStatement := "select fundings.id, users.nick_name, fundings.main_img_url, fundings.end_date - fundings.start_date as d_day, " +
		"fundings.name, fundings.sales_amount, fundings.start_date, fundings.end_date, fundings.price, " +
		"fundings.target_amount, artists.name, fundings.sales_amount / fundings.target_amount * 100 as achievement_rate from users " +
		"inner join fundings on fundings.seller_id = users.id inner join artists on fundings.artist_id = artists.id " +
		"where fundings.id = ? and fundings.deleted_at is null"

	row := configs.DB.Debug().Raw(sqlStatement, fundID).Row()
	err := row.Scan(&body.ID, &body.NickName, &titleImg, &body.Dday, &body.Name, &body.SalesAmount, &body.StartDate, &body.EndDate,
		&body.Price, &body.TargetAmount, &body.ArtistName, &body.AchievementRate)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("not found")
	}
	body.FundingImgUrls = append(body.FundingImgUrls, titleImg)

	sqlStatement = "select url, is_title from funding_imgs " +
		"where funding_imgs.funding_id = ? and funding_imgs.deleted_at is null " +
		"order by funding_imgs.order"
	rows, _ := configs.DB.Raw(sqlStatement, fundID).Rows()
	defer rows.Close()
	for rows.Next() {
		var url string
		var isTitle bool
		rows.Scan(&url, &isTitle)
		if isTitle {
			body.FundingImgUrls = append(body.FundingImgUrls, url)
		} else {
			body.DetailedImgUrl = url
		}
	}
	return &body, nil
}

type queryString struct {
	ArtistID uint `form:"artist-id" binding:"required"`
}

type fundingListResBody struct {
	ID              uint    `json:"id"`
	NickName        string  `json:"sellerName"`
	Name            string  `json:"fundingName"`
	MainImgUrl      string  `json:"mainImgUrl"`
	DDay            int     `json:"dDay"`
	AchievementRate float64 `json:"achievementRate"`
}

func GetFundingList(c *gin.Context) {
	queryString := queryString{}

	err := c.ShouldBindQuery(&queryString)
	if err != nil || !isArtistExist(queryString.ArtistID) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "no artist",
		})
	} else {
		body := setFundingListBody(queryString.ArtistID)
		c.JSON(http.StatusOK, gin.H{
			"data": body,
		})
	}

}

func setFundingListBody(artistID uint) []fundingListResBody {
	body := []fundingListResBody{}

	configs.DB.Debug().Table("fundings").Joins("inner join users on fundings.seller_id = users.id").
		Select("fundings.id, users.nick_name, fundings.name, fundings.main_img_url, "+
			"fundings.end_date - fundings.start_date as d_day, fundings.sales_amount / fundings.target_amount * 100 as achievement_rate").
		Where("fundings.artist_id = ? and fundings.deleted_at is null", artistID).Order("d_day").
		Scan(&body)
	return body
}

func isArtistExist(artistID uint) bool {
	artists := models.Artist{}

	err := configs.DB.First(&artists, artistID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else {
		return true
	}

}

type fundingstringBody struct {
	ID             uint     `json:"id"`
	NickName        string   `json:"sellerName"`
	Name            string   `json:"fundName"`
	Price           uint     `json:"price"`
	TargetAmount    uint     `json:"targetAmount"`
	SalesAmount     uint     `json:"salesAmount"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	ArtistName      string   `json:"artistName"`
	AchievementRate float64  `json:"achievementRate"` //salesAmount / Price
	Dday            uint     `json:"dDay"`
	FundingImgUrls  []string `json:"fundingImgUrls"`
	DetailedImgUrl  string   `json:"detailedImgUrl"`
}

func CreateFund(c *gin.Context) {
	var funding models.Funding
	c.BindJSON(&funding)

	configs.DB.Create(&funding)
	c.JSON(http.StatusOK, funding)
}
