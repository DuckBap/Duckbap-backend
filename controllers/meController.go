package controllers

import (
	"fmt"
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type profile struct {
	NickName string
	ArtistID uint
	ImgUrl string
}

type ranking struct {
	UserID uint
	ArtistID uint
	SellTotal uint
	BuyTotal uint
	Total uint
}

type buy struct {
	MainImgUrl string
	Name string
	CreatedAt time.Time
	BuyerID uint

}

type sell struct {
	MainImgUrl string
	Name string
	CreatedAt time.Time
	SellerID uint

}

type bookmark struct {
	ArtistID uint
	UserID uint
	Level int
}


/*
func getRanking(ranking *ranking, user models.User) {
	configs.DB.Model(&models.Funding{}).
		Select("IFNULL(sum(sales_amount * price), 0)").
		Where("artist_id = ? AND seller_id = ?", user.Artist.ID, user.ID).
		Scan(&ranking.SellTotal)

	configs.DB.Model(&models.Receipt{}).
		Select("IFNULL(sum(receipts.amount * fundings.price), 0)").
		Joins("join fundings  on fundings.id = receipts.funding_id").
		Where("artist_id = ? AND receipts.buyer_id = ?", user.Artist.ID, user.ID).
		Scan(&ranking.BuyTotal)

	ranking.Total = ranking.SellTotal + ranking.BuyTotal
	ranking.UserID = user.ID
	ranking.ArtistID = user.Artist.ID
}
*/
func getProfile(user models.User) profile {
	var profile profile
	profile.ImgUrl = user.Artist.ImgUrl
	profile.NickName = user.NickName
	profile.ArtistID = user.Artist.ID

	/*
		configs.DB.Model(&user).
			Select("nick_name, artists.img_url, artists.id").
			Joins("join artists on users.favorite_artist = artists.id").
			Where("users.id = ?", userid).
			Scan(&profile)
	*/
	return profile
}

func getSell(user models.User) []sell {
	var sell []sell
	configs.DB.Table("fundings").
		Select("name, main_img_url, receipts.created_at, receipts.seller_id").
		Joins("join receipts on fundings.id = receipts.funding_id").
		Where("receipts.seller_id = ?", user.ID).
		Order("receipts.created_at desc").
		Limit(2).
		Scan(&sell)
	return sell
}

func getBuy(user models.User) []buy {
	var buy []buy
	configs.DB.Table("fundings").
		Select("name, main_img_url, receipts.created_at, receipts.buyer_id").
		Joins("join receipts on fundings.id = receipts.funding_id").
		Where("receipts.buyer_id = ?", user.ID).
		Order("receipts.created_at desc").
		Limit(2).
		Scan(&buy)
	return buy
}

func getRankingList(user models.User) []ranking {
	var rankingList []ranking
	var temp []ranking
	configs.DB.Table("Users").
		Select("favorite_artist as artist_id, id as user_id").
		Where("favorite_artist = ?", user.Artist.ID).
		Find(&rankingList)

	configs.DB.Table("Bookmarks").
		Select("artist_id, user_id").
		Where("artist_id = ?", user.Artist.ID).
		Find(&temp)

	for i := range temp {
		rankingList = append(rankingList, temp[i])
	}

	for i := range rankingList {
		configs.DB.Model(&models.Funding{}).
			Select("IFNULL(sum(sales_amount * price), 0)").
			Where("artist_id = ? AND seller_id = ?", user.Artist.ID, rankingList[i].UserID).
			Scan(&rankingList[i].SellTotal)

		configs.DB.Model(&models.Receipt{}).
			Select("IFNULL(sum(receipts.amount * fundings.price), 0)").
			Joins("join fundings on fundings.id = receipts.funding_id").
			Where("artist_id = ? AND receipts.buyer_id = ?", user.Artist.ID, rankingList[i].UserID).
			Scan(&rankingList[i].BuyTotal)
		rankingList[i].Total = rankingList[i].SellTotal + rankingList[i].BuyTotal
	}
	sort.Slice(rankingList, func(i, j int) bool {
		return rankingList[i].Total > rankingList[j].Total
	})
	return rankingList
}

func associate(user *models.User) error {
	/*
	if err := configs.DB.Model(&user).Association("Fundings").Find(&user.Fundings); err != nil {
		return err
	}
	if err := configs.DB.Model(&user).Association("Receipts").Find(&user.Receipts); err != nil {
		return err
	}*/
	if err := configs.DB.Model(&user).Association("Artist").Find(&user.Artist); err != nil {
		return err
	}
	return nil
}

func getLevel(rankingList []ranking, user models.User) int {
	var rank int
	var level int

	for i, idx := range rankingList {
		if idx.UserID == user.ID {
			rank = i + 1
			break
		}
	}
	if myrank := float32(rank) / float32(len(rankingList)) * 100 ; myrank <= 3 {
		level = 1
	} else if myrank <= 10 {
		level = 2
	} else if myrank <= 20 {
		level = 3
	} else if myrank <= 50 {
		level = 4
	} else {
		level = 5
	}
	return level
}

func getBookmark(user models.User) []bookmark {
	var bookmark []bookmark
	configs.DB.Where("user_id", user.ID).Find(&bookmark)

	for idx, _ := range bookmark {
		var bookmarkRanking []ranking
		var temp []ranking

		configs.DB.Table("Users").
			Select("favorite_artist as artist_id, id as user_id").
			Where("favorite_artist = ?", bookmark[idx].ArtistID).
			Find(&bookmarkRanking)
		configs.DB.Table("Bookmarks").
			Select("artist_id, user_id").
			Where("artist_id = ?", bookmark[idx].ArtistID).
			Find(&temp)
		for i := range temp {
			bookmarkRanking = append(bookmarkRanking, temp[i])
		}
		for i := range bookmarkRanking {
			configs.DB.Model(&models.Funding{}).
				Select("IFNULL(sum(sales_amount * price), 0)").
				Where("artist_id = ? AND seller_id = ?", bookmark[idx].ArtistID, bookmarkRanking[i].UserID).
				Scan(&bookmarkRanking[i].SellTotal)

			configs.DB.Model(&models.Receipt{}).
				Select("IFNULL(sum(receipts.amount * fundings.price), 0)").
				Joins("join fundings on fundings.id = receipts.funding_id").
				Where("artist_id = ? AND receipts.buyer_id = ?", bookmark[idx].ArtistID, bookmarkRanking[i].UserID).
				Scan(&bookmarkRanking[i].BuyTotal)
			bookmarkRanking[i].Total = bookmarkRanking[i].SellTotal + bookmarkRanking[i].BuyTotal
		}
		sort.Slice(bookmarkRanking, func(i, j int) bool {
			return bookmarkRanking[i].Total > bookmarkRanking[j].Total
		})
		bookmark[idx].Level = getLevel(bookmarkRanking, user)
	}
	sort.Slice(bookmark, func(i, j int) bool {
		return bookmark[i].Level < bookmark[j].Level
	})
	return bookmark
}

func GetMe(c *gin.Context) {
	var user models.User
	userid, err := strconv.Atoi(c.PostForm("userid"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := configs.DB.First(&user, userid).Error; err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := associate(&user) ; err != nil {
		fmt.Println(err.Error())
		return
	}

	//getRanking(&ranking, user)

	profile := getProfile(user)
	buy := getBuy(user)
	sell := getSell(user)
	rankingList := getRankingList(user)
	level := getLevel(rankingList, user)
	bookmark := getBookmark(user)


	c.JSON(http.StatusOK, gin.H{
		"profile":  profile,
		"buy": buy,
		"sell": sell,
		"level": level,
		"bookmark": bookmark,
	})
}
