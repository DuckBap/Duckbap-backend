package controllers

import (
	"fmt"
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/permissions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type	InputUserData struct {
	UserName string
	Password1 string
	Password2 string
	Email string
	NickName string
	FavoriteArtist uint
}

func hash(pwd string) string {
	digest, _ := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(digest)
}

func	SignUp (c *gin.Context) {
	var	user		models.User
	var	inputData	InputUserData
	var	errorPoint	string
	var httpCode	int
	var	checker		bool

	err := c.Bind(&inputData)
	if err != nil {
		c.JSON(400, err)
	}
	errorPoint, httpCode, checker = permissions.IsEmpty(&inputData)
	if checker {
		fmt.Println("Impossible")
		c.JSON(httpCode, errorPoint)
		return
	}
	user.UserName = inputData.UserName
	user.Password = inputData.Password1
	user.NickName = inputData.NickName
	user.FavoriteArtist = inputData.FavoriteArtist
	errorPoint, httpCode, checker = permissions.IsExist(&user)
	if checker {
		fmt.Println("Impossible")
		c.JSON(httpCode, errorPoint)
		return
	}
	password := hash(user.Password)
	user.Password = password
	tx := configs.DB.Create(&user)
	if tx.Error != nil {
		errorPoint, httpCode = permissions.FindErrorPoint(tx.Error)
		c.JSON(httpCode, errorPoint)
		return
	}
	c.JSON(httpCode, user)
}

/* url : Get /sign-up
** 아티스트의 목록을 보내줘서 보여줘야 한다.
** 이유 : 회원 가입시 필수로 최애 아티스트를 선택해야 되기 때문이다.
 */
func	ShowArtists (c *gin.Context) {
	var artist		models.Artist
	var	artists		[]models.Artist
	var	errorPoint	string
	var httpCode	int

	tx := configs.DB.Model(&artist).Select("id, name").Find(&artists)
	if tx.Error != nil {
		errorPoint = tx.Error.Error()
		httpCode = http.StatusNotFound
	} else {
		httpCode = http.StatusOK
	}
	if errorPoint != "" {
		c.JSON(httpCode, errorPoint)
		return
	}
	c.JSON(httpCode, artists)
}