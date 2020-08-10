package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
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

type	itemList	struct {
	MainImgUrl	string	`gorm:"varchar(255); unique; not null"`
	Name		string	`gorm:"varchar(150);not null;"`
	AchievementRate	float64	`json:"achievementRate"`
}

func BannerSelect (c *gin.Context) {
	var fundings []bannerFunding

	configs.DB.Table("fundings").Order("sales_amount desc").Limit(5).Find(&fundings)
	c.JSON(200, fundings)
}

func ListSelect (c *gin.Context) {
	var bookmark []bookmarks
	var favorite []bookmarks
	var fundings []listFunding
	var tmp []listFunding

	user, _ := c.Get("user")
	id := user.(*models.User).ID
	configs.DB.Table("bookmarks").Where("user_id = ?", id).Order("artist_id").Find(&bookmark)
	configs.DB.Table("users").Select("favorite_artist").Where("user_id = ?", id).Find(&favorite)
	bookmark = append(bookmark, favorite...)
	limit := int(math.Ceil(8.0/float64(len(bookmark))))
	for i:=0;i<len(bookmark);i++ {
		configs.DB.Table("fundings").Where("artist_id = ?", bookmark[i].ArtistID).Order("sales_amount desc").Limit(limit).Find(&tmp)
		fundings = append(fundings, tmp...)
	}
	sort.Slice(fundings, func (i, j int) bool {
		return fundings[i].SalesAmount > fundings[j].SalesAmount
	})
	if len(fundings) < 8 {
		dup := setDuplicates(bookmark)
		configs.DB.Table("fundings").Not("artist_id", dup).Order("sales_amount desc").Limit(8 - len(fundings)).Find(&tmp)
		fundings = append(fundings, tmp...)
	}
	c.JSON(200, fundings)
}

func	NotloginListSelect(c *gin.Context) {
	var list []itemList

	configs.DB.Raw("select main_img_url, name, (@achievement_rate:=100 * sales_amount/target_amount)achievement_rate from (select f.*, (case @vartist when f.artist_id then @rownum:=@rownum+1 else @rownum:=1 end)rnum, (@vartist:=f.artist_id)vartist from(select * from fundings order by artist_id, sales_amount desc)f, (select @vartist:='',@rownum:=0 from dual)b)e where rnum <= 2 order by sales_amount desc limit 8").Scan(&list)
	c.JSON(http.StatusOK, list)
}

func setDuplicates (bookmark []bookmarks) []uint {
	var dup []uint

	for i:=0;i<len(bookmark);i++ {
		dup[i] = bookmark[i].ArtistID
	}
	return dup
}
