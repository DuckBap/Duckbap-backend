package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type listFunding struct {
	ID              uint
	Name            string
	TargetAmount    uint `json:"-"`
	MainImgUrl      string
	ArtistID        uint    `json:"-"`
	SalesAmount     uint    `json:"-"`
	AchievementRate float64 `json:"achievementRate"`
}

type bannerFunding struct {
	ID         uint
	Name       string
	MainImgUrl string
	ArtistID   uint
}

type bookmarks struct {
	ArtistID uint
}

type itemList struct {
	AchievementRate float64       `json:"achievementRate"`
	Dday            time.Duration `json:"dDay"`
	Name            string        `gorm:"varchar(150);not null;"`
	ID              uint
	MainImgUrl      string `gorm:"varchar(255); unique; not null"`
	NickName        string
}

// @Summary 메인 배너에서 보여줄 펀딩 리스트
// @Description <br>펀딩 리스트를 반환합니다.<br>
// @Accept  json
// @Produce  json
// @Router /fundings/banner [get]
// @Success 200 {array} bannerFunding
func BannerSelect(c *gin.Context) {
	var fundings []bannerFunding

	configs.DB.Table("fundings").Order("sales_amount desc").Limit(5).Find(&fundings)
	c.JSON(http.StatusOK, gin.H{
		"data": fundings,
	})
}

// @Summary 메인에서 보여줄 펀딩 리스트
// @Description ## 메인 페이지에서 보여줄 펀딩 리스트를 반환합니다.
// @Description <br>
// @Description ## 로그인이 되어있을 경우
// @Description 1. 즐겨찾기에 저장되어있는 아이돌, 최애 아이돌과 관련된 펀딩들, 판매량이 가장 높은 펀딩들이 포함됩니다.
// @Description <br>
// @Description 2. 펀딩 8개가 들어있는 리스트가 반환됩니다.<br>
// @Description <br>
// @Description ## 로그인이 되어있지 않을 경우
// @Description 1. 판매량이 높은 펀딩 8개가 포함된 리스트가 반환됩니다.
// @Accept  json
// @Produce  json
// @Router /fundings/main [get]
// @Success 200 {array} listFunding
func ListSelect(c *gin.Context, id uint) {
	var bookmark []bookmarks
	type Test struct {
		FavoriteArtist uint
	}
	var favorite Test
	var fundings []listFunding
	var tmp []listFunding
	var temp bookmarks

	configs.DB.Table("bookmarks").Where("user_id = ?", id).Order("artist_id").Find(&bookmark)
	configs.DB.Table("users").Select("favorite_artist").Where("id = ?", id).Find(&favorite)
	temp.ArtistID = favorite.FavoriteArtist
	bookmark = append(bookmark, temp)
	limit := int(math.Ceil(8.0 / float64(len(bookmark))))
	for i := 0; i < len(bookmark); i++ {
		configs.DB.Table("fundings").Where("artist_id = ?", bookmark[i].ArtistID).Order("sales_amount desc").Limit(limit).Find(&tmp)
		fundings = append(fundings, tmp...)
	}
	sort.Slice(fundings, func(i, j int) bool {
		return fundings[i].SalesAmount > fundings[j].SalesAmount
	})
	if len(fundings) < 8 {
		dup := setDuplicates(bookmark)
		configs.DB.Table("fundings").Where("artist_id Not In (?)", dup).Order("sales_amount desc").Limit(8 - len(fundings)).Find(&tmp)
		fundings = append(fundings, tmp...)
	}
	for i, item := range fundings {
		tmp_rate := float64(item.SalesAmount) / float64(item.TargetAmount) * 10000
		int_rate := int(tmp_rate)
		fundings[i].AchievementRate = float64(int_rate) / 100
	}
	c.JSON(http.StatusOK, gin.H{
		"data": fundings,
	})
}

type test struct {
	Dday time.Duration
}

func NotloginListSelect(c *gin.Context) {
	temp, _ := c.GetQuery("limit")
	limit, _ := strconv.Atoi(temp)
	var list []itemList

	configs.DB.Raw("select (@achievement_rate:=truncate(100 * sales_amount/target_amount,2))achievement_rate, " +
						"(@Dday:=datediff(end_date,now()))Dday, e.id, name, main_img_url, users.nick_name "+
						"from ("+
							"select f.*,"+
								"(case @vartist when f.artist_id then @rownum:=@rownum+1 else @rownum:=1 end)rnum, "+
								"(@vartist:=f.artist_id)vartist " +
							"from(select * from fundings order by artist_id, sales_amount desc)f,"+
								"(select @vartist:='',@rownum:=0 from dual)b" +
						")e "+
						"left join users on seller_id = users.id " +
						"where rnum <= 2 order by sales_amount desc limit ?, ?", limit, 8).Scan(&list)
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func setDuplicates(bookmark []bookmarks) []uint {
	var dup []uint

	for i := 0; i < len(bookmark); i++ {
		dup = append(dup, bookmark[i].ArtistID)
	}
	return dup
}
