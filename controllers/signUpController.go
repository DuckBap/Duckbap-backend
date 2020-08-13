package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/permissions"
	//"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type InputUserData struct {
	UserName       string `form:"userName" json:"userName"`
	Password1      string `form:"password1" json:"password1"`
	Password2      string `form:"password2" json:"password2"`
	Email          string `form:"email" json:"email"`
	NickName       string `form:"nickName" json:"nickName"`
	FavoriteArtist uint   `form:"favoriteArtist" json:"favoriteArtist"`
}

func inputDataToUser(user *models.User, inputData InputUserData) {
	(*user).UserName = inputData.UserName
	(*user).Password = inputData.Password1
	(*user).NickName = inputData.NickName
	(*user).Email = inputData.Email
	(*user).FavoriteArtist = inputData.FavoriteArtist
}

func hash(pwd string) string {
	digest, _ := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(digest)
}

func SignUp(c *gin.Context) {
	var user 		models.User
	var inputData	InputUserData
	var httpCode	int
	var checker		bool

	err := c.ShouldBind(&inputData)
	errorMap := make(map[string]string)
	if err != nil {
		errorPoint := permissions.AnalyzeErrorMessage(err.Error())
		errorMap[errorPoint] = "잘못된 값입니다."
		c.JSON(http.StatusBadRequest, gin.H {
			"err": errorMap,
		})
		c.Abort()
		return
	}
	errorMap, httpCode, checker = permissions.IsEmpty(&inputData)
	if checker {
		c.JSON(httpCode, gin.H {
			"err": errorMap,
		})
		c.Abort()
		return
	}
	inputDataToUser(&user, inputData)
	errorMap, httpCode, checker = permissions.IsExist(&user)
	if checker {
		c.JSON(httpCode, gin.H {
			"err": errorMap,
		})
		c.Abort()
		return
	}
	password := hash(user.Password)
	user.Password = password
	tx := configs.DB.Create(&user)
	if tx.Error != nil {
		errorMap, httpCode = permissions.FindErrorPoint(tx.Error)
		c.JSON(httpCode, gin.H {
			"err": errorMap,
		})
		c.Abort()
		return
	}
	c.Set("user", &user)
}
