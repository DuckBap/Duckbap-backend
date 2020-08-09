package controller

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
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

type FundingResBody struct{
	NickName			string		`json:"sellerName"`
	Name				string		`json:"fundName"`
	Price				uint		`json:"price"`
	TargetAmount		uint		`json:"targetAmount"`
	SalesAmount			uint		`json:"salesAmount"`
	StartDate			time.Time	`json:"startDate"`
	EndDate				time.Time	`json:"endDate"`
	ArtistName			string		`json:"artistName"`
	AchievementRate		float32		`json:"achievementRate"`	//salesAmount / Price
	Dday				uint		`json:"dDay"`
	FundingImgUrls		[]string	`json:"fundingImgUrls"`
	DetailedImgUrl		string		`json:"detailedImgUrl"`
}


//test
func CreateFunding(c *gin.Context) {
	fund := models.Funding{
		SellerID: 1,
		Name: "아이유 굿즈3",
		Price: 4000,
		TargetAmount: 50000,
		MainImgUrl: "이미지3",
		ArtistID: 1,
		StartDate: time.Now(),
		EndDate: time.Now().Add(24 * time.Hour),
	}
	configs.DB.Create(&fund)
	c.JSON(http.StatusOK, gin.H{
		"msg": "create funding",
		"funding": fund,
	})
}

func GetFunding(c *gin.Context) {
	fundID := c.Param("fund_id")

	body, err := SetFundingBody(fundID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"mgs": "no funding",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "get funding",
		"funding": body,
	})
}

func SetFundingBody(fundID string) (*FundingResBody, error){
	body := FundingResBody{}
	var titleImg string

	sqlStatement := "select users.nick_name, fundings.main_img_url, fundings.end_date - fundings.start_date as d_day, fundings.name, fundings.sales_amount, fundings.start_date, fundings.end_date, fundings.price, fundings.target_amount, artists.name,	fundings.sales_amount / fundings.target_amount as achievement_rate from users inner join fundings on fundings.seller_id = users.id inner join artists on fundings.artist_id = artists.id where fundings.id = ? and fundings.deleted_at is null"
	row := configs.DB.Debug().Raw(sqlStatement, fundID).Row()
	row.Scan(&body.NickName, &titleImg, &body.Dday, &body.Name, &body.SalesAmount, &body.StartDate, &body.EndDate, &body.Price, &body.TargetAmount, &body.ArtistName, &body.AchievementRate)
	body.FundingImgUrls = append(body.FundingImgUrls, titleImg)

	sqlStatement = "select url, is_title from funding_imgs where funding_imgs.funding_id = ? and funding_imgs.deleted_at is null order by funding_imgs.order"
	rows, err := configs.DB.Raw(sqlStatement, fundID).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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

type QueryString struct {
	ArtistID string	`form:"artist-id" binding:"required"`
}

type FundingListResBody struct {
	NickName			string		`json:"sellerName"`
	Name				string		`json:"fundingName"`
	MainImgUrl			string		`json:"mainImgUrl"`
	DDay				int			`json:"dDay"`
	AchievementRate		float64		`json:"achievementRate"`
}

func GetFundingList(c *gin.Context) {
	queryString := QueryString{}
	c.BindQuery(&queryString)

	body := SetFundingListBody(queryString.ArtistID)
	c.JSON(http.StatusOK, gin.H{
		"fundList": body,
	})
}

func SetFundingListBody(artistID string) []FundingListResBody{
	body := []FundingListResBody{}

	configs.DB.Table("fundings").Joins("inner join users on fundings.seller_id = users.id join artists").
		Select("users.nick_name, fundings.name, fundings.main_img_url, fundings.end_date - fundings.start_date as d_day, fundings.sales_amount / fundings.target_amount as achievement_rate").
		Where("artists.id = ? and fundings.deleted_at is null", artistID).Order("d_day").
		Scan(&body)
	return body
}

func BuyFunding(c *gin.Context) {

}