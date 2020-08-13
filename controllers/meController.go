package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
	"sort"
	"time"
)

type data struct {
	ID             uint      `json:"id"`
	Nickname       string    `json:"nickName"`
	FavoriteArtist artist    `json:"favoriteArtist"`
	Buy            []receipt `json:"buy"`
	Sell           []receipt `json:"sell"`
	Bookmark       []artist  `json:"bookmark"`
}

type artist struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	ImgUrl string `json:"imgUrl"`
	Level  int    `json:"level"`
}

type ranking struct {
	UserID    uint
	ArtistID  uint
	SellTotal uint
	BuyTotal  uint
	Total     uint
}

type funding struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	MainImgUrl string `json:"mainImgUrl"`
}

type receipt struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FundingID uint      `json:"fundingId"`
	Funding   funding   `json:"funding"`
}

func getArtist(user models.User, rankingList []ranking) artist {
	var artist artist

	artist.ID = user.Artist.ID
	artist.ImgUrl = user.Artist.ImgUrl
	artist.Name = user.Artist.Name
	artist.Level = getLevel(rankingList, user)
	return artist
}

func getSell(user models.User) ([]receipt, error) {
	var sell []receipt

	err := configs.DB.Table("receipts").
		Select("receipts.id, fundings.id as funding_id").
		Joins("join fundings on fundings.id = receipts.funding_id").
		Where("receipts.seller_id = ?", user.ID).
		Order("receipts.created_at desc").
		Limit(2).
		Preload(clause.Associations).
		Find(&sell).
		Error
	if err != nil {
		return nil, err
	}
	return sell, nil
}

func getBuy(user models.User) ([]receipt, error) {
	var buy []receipt

	err := configs.DB.Table("receipts").
		Select("receipts.id, fundings.id as funding_id").
		Joins("join fundings on fundings.id = receipts.funding_id").
		Where("receipts.buyer_id = ?", user.ID).
		Order("receipts.created_at desc").
		Limit(2).
		Preload(clause.Associations).
		Find(&buy).
		Error
	if err != nil {
		return nil, err
	}
	return buy, nil
}

func getRankingList(user models.User) ([]ranking, error) {
	var rankingList []ranking
	var temp []ranking
	err := configs.DB.Table("Users").
		Select("favorite_artist as artist_id, id as user_id").
		Where("favorite_artist = ?", user.Artist.ID).
		Find(&rankingList).
		Error
	if err != nil {
		return nil, err
	}
	err = configs.DB.Table("Bookmarks").
		Select("artist_id, user_id").
		Where("artist_id = ?", user.Artist.ID).
		Find(&temp).
		Error
	if err != nil {
		return nil, err
	}

	for i := range temp {
		rankingList = append(rankingList, temp[i])
	}

	for i := range rankingList {
		err := configs.DB.Model(&models.Funding{}).
			Select("IFNULL(sum(sales_amount * price), 0)").
			Where("artist_id = ? AND seller_id = ?", user.Artist.ID, rankingList[i].UserID).
			Scan(&rankingList[i].SellTotal).
			Error
		if err != nil {
			return nil, err
		}

		err = configs.DB.Model(&models.Receipt{}).
			Select("IFNULL(sum(receipts.amount * fundings.price), 0)").
			Joins("join fundings on fundings.id = receipts.funding_id").
			Where("artist_id = ? AND receipts.buyer_id = ?", user.Artist.ID, rankingList[i].UserID).
			Scan(&rankingList[i].BuyTotal).
			Error
		if err != nil {
			return nil, err
		}
		rankingList[i].Total = rankingList[i].SellTotal + rankingList[i].BuyTotal
	}
	sort.Slice(rankingList, func(i, j int) bool {
		return rankingList[i].Total > rankingList[j].Total
	})
	return rankingList, nil
}

func associate(user *models.User) error {
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
	if myrank := float32(rank) / float32(len(rankingList)) * 100; myrank <= 3 {
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

func getBookmark(user models.User) ([]artist, error) {
	var bookmark []artist

	err := configs.DB.Table("bookmarks").
		Select("bookmarks.artist_id as id, artists.name, artists.img_url").
		Joins("join artists on bookmarks.artist_id = artists.id").
		Where("user_id", user.ID).
		Find(&bookmark).
		Error
	if err != nil {
		return nil, err
	}

	for idx, _ := range bookmark {
		var bookmarkRanking []ranking
		var temp []ranking

		err := configs.DB.Table("Users").
			Select("favorite_artist as artist_id, id as user_id").
			Where("favorite_artist = ?", bookmark[idx].ID).
			Find(&bookmarkRanking).
			Error
		if err != nil {
			return nil, err
		}

		err = configs.DB.Table("Bookmarks").
			Select("artist_id, user_id").
			Where("artist_id = ?", bookmark[idx].ID).
			Find(&temp).
			Error
		if err != nil {
			return nil, err
		}

		for i := range temp {
			bookmarkRanking = append(bookmarkRanking, temp[i])
		}
		for i := range bookmarkRanking {
			err := configs.DB.Model(&models.Funding{}).
				Select("IFNULL(sum(sales_amount * price), 0)").
				Where("artist_id = ? AND seller_id = ?", bookmark[idx].ID, bookmarkRanking[i].UserID).
				Scan(&bookmarkRanking[i].SellTotal).
				Error
			if err != nil {
				return nil, err
			}

			err = configs.DB.Model(&models.Receipt{}).
				Select("IFNULL(sum(receipts.amount * fundings.price), 0)").
				Joins("join fundings on fundings.id = receipts.funding_id").
				Where("artist_id = ? AND receipts.buyer_id = ?", bookmark[idx].ID, bookmarkRanking[i].UserID).
				Scan(&bookmarkRanking[i].BuyTotal).
				Error
			if err != nil {
				return nil, err
			}
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
	return bookmark, nil
}

func GetMe(c *gin.Context) {
	var user models.User
	var data data

	loginUser, _ := c.Get("user") //로그인한 유저 인터페이스 가져오기
	if err := configs.DB.First(&user, loginUser.(*models.User).ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	if err := associate(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	data.ID = user.ID
	data.Nickname = user.NickName

	if sell, err := getSell(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	} else {
		data.Sell = sell
	}

	if buy, err := getBuy(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	} else {
		data.Buy = buy
	}

	rankingList, err := getRankingList(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	data.FavoriteArtist = getArtist(user, rankingList)
	if bookmark, err := getBookmark(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	} else {
		data.Bookmark = bookmark
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

type bookmark struct {
	UserID   uint `json:"user_id"`
	ArtistID uint `json:"artist_id"`
}

func CreateBookmark(c *gin.Context) {
	var bookmarkrecord bookmark
	c.BindJSON(&bookmarkrecord)

	configs.DB.Create(&bookmarkrecord)
	c.JSON(http.StatusOK, bookmarkrecord)
}
