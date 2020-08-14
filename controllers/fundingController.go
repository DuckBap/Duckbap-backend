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

// @Summary 펀딩 상세정보
// @Description <br>펀딩 상세정보를 반환합니다.
// @Description <br>
// @Description sellerName : 판매자의 닉네임<br>
// @Description fundName : 펀드 이름<br>
// @Description price : 하나를 구매할 때의 가격<br>
// @Description targetAmount : 판매 목표량<br>
// @Description salesAmount : 현재까지의 판매량<br>
// @Description startDate : 펀딩 시작 일<br>
// @Description endDate: 펀딩 마감 일<br>
// @Description artistName : 펀딩과 관련된 연예인 이름<br>
// @Description achievementRate : 펀딩 달성률 (판매량 / 목표량)<br>
// @Description dDay : 펀딩 마감일까지 남은 날짜<br>
// @Description fundingImgUrls : 펀딩 상품의 이미지 주소들<br>
// @Description detailedImgUrl : 펀딩 상세정보 이미지<br>
// @Param fund_id path integer true "fund_id"
// @Accept  json
// @Produce  json
// @Router /fundings/{fund_id} [get]
// @Success 200 {object} fundingResBody
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

// @Summary 아티스트와 관련 펀딩 리스트
// @Description <br>아티스트와 관련된 펀딩 리스트를 반환합니다.
// @Description <br>
// @Param artist-id query integer true "."
// @Accept  json
// @Produce  json
// @Router /fundings [get]
// @Success 200 {array} fundingListResBody
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
	ID              uint     `json:"id"`
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

func CreateFundingImg(c *gin.Context) {
	var img models.FundingImg
	c.BindJSON(&img)

	configs.DB.Create(&img)
	c.JSON(http.StatusOK, img)
}
