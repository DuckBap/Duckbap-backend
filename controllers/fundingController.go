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
	SellerName 			string
	FundName 			string
	Price				uint
	TargetAmount 		uint
	SalesAmount 		uint
	StartDate 			time.Time
	EndDate 			time.Time
	ArtistName 			string
	AchievementRate		float32			//salesAmount / Price
	Dday				uint
	FundingImgUrls		[]string
	DetailedImgUrl		string
}

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
	body := FundingResBody{}
	fundID := c.Param("fund_id")
	SetFundingBody(&body, fundID)
	c.JSON(http.StatusOK, gin.H{
		"msg": "get funding",
		"funding": body,
	})
}

func SetFundingBody(body *FundingResBody, fundID string) {
	var titleImg string

	row := configs.DB.Raw("select users.nick_name, fundings.main_img_url, fundings.end_date - fundings.start_date as d_day, fundings.name, fundings.sales_amount, fundings.start_date, fundings.end_date, fundings.price, fundings.target_amount, artists.name,	fundings.sales_amount / fundings.target_amount as achievement_rate from users inner join fundings on fundings.seller_id = users.id inner join artists on fundings.artist_id = artists.id where fundings.id = ? and fundings.deleted_at is null", fundID).Row()
	row.Scan(&body.SellerName, &titleImg, &body.Dday, &body.FundName, &body.SalesAmount, &body.StartDate, &body.EndDate, &body.Price, &body.TargetAmount, &body.ArtistName, &body.AchievementRate)
	body.FundingImgUrls = append(body.FundingImgUrls, titleImg)

	rows, _ := configs.DB.Raw("select url, is_title from funding_imgs inner join fundings on funding_imgs.funding_id = ? where funding_imgs.deleted_at is null order by funding_imgs.order", fundID).Rows()
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
}

type QueryString struct {
	ArtistID string	`form:"artist-id" binding:"required"`
}

type FundingListResBody struct {
	SellerName      string
	FundingName     string
	MainImgUrl      string
	DDay            int
	AchievementRate float32
}

func GetFundingList(c *gin.Context) {
	queryString := QueryString{}
	c.BindQuery(&queryString)

	body := SetFundingListBody(queryString.ArtistID)

	c.JSON(http.StatusOK, gin.H{
		"fund_list": body,
	})
}

func SetFundingListBody(artistID string) *[]FundingListResBody{
	body := []FundingListResBody{}

	rows, _ := configs.DB.Raw("select users.user_name, fundings.name, fundings.main_img_url, fundings.end_date - fundings.start_date as d_day, fundings.sales_amount / fundings.target_amount as acheivement_rate from fundings inner join users on fundings.seller_id = users.id join artists where artists.id = ? and fundings.deleted_at is null order by d_day", artistID).Rows()
	defer rows.Close()
	for rows.Next() {
		var seller 				string
		var fundingName 		string
		var imgUrl				string
		var dDay				int
		var achievementRate		float32

		rows.Scan(&seller, &fundingName, &imgUrl, &dDay, &achievementRate)
		body = append(body, FundingListResBody{
			SellerName: seller,
			FundingName: fundingName,
			MainImgUrl: imgUrl,
			DDay: dDay,
			AchievementRate: achievementRate,
		})
	}
	return &body
}

func BuyFunding(c *gin.Context) {

}