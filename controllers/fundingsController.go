package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"sort"
)

type listFunding struct {
	ID uint
	//SellerID uint
	Name string
	//Price uint
	TargetAmount    uint `json:"-"`
	MainImgUrl      string
	ArtistID        uint `json:"-"`
	SalesAmount     uint `json:"-"`
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
	ID              uint
	MainImgUrl      string  `gorm:"varchar(255); unique; not null"`
	Name            string  `gorm:"varchar(150);not null;"`
	AchievementRate float64 `json:"achievementRate"`
}

func BannerSelect(c *gin.Context) {
	var fundings []bannerFunding

	configs.DB.Table("fundings").Order("sales_amount desc").Limit(5).Find(&fundings)
	c.JSON(200, fundings)
}

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
	c.JSON(200, fundings)
}

func NotloginListSelect(c *gin.Context) {
	var list []itemList

	configs.DB.Raw("select id, main_img_url, name, (@achievement_rate:=truncate(100 * sales_amount/target_amount,2))achievement_rate " +
						"from (" +
								"select f.*," +
									"(case @vartist when f.artist_id then @rownum:=@rownum+1 else @rownum:=1 end)rnum," +
									"(@vartist:=f.artist_id)vartist " +
								"from(" +
									"select * " +
									"from fundings order by artist_id, sales_amount desc)f," +
									"(select @vartist:='',@rownum:=0 from dual)b)e " +
									"where rnum <= 2 order by sales_amount desc limit 8").Scan(&list)
	c.JSON(http.StatusOK, list)
}

func setDuplicates(bookmark []bookmarks) []uint {
	var dup []uint

	for i := 0; i < len(bookmark); i++ {
		dup = append(dup, bookmark[i].ArtistID)
	}
	return dup
}
