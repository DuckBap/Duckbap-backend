package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/permissions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type InputUserData struct {
	UserName       string `form:"userName" json:"userName"`
	Password1      string `form:"password1" json:"password1"`
	Password2      string `form:"password2" json:"password2"`
	Email          string `form:"email" json:"email"`
	NickName       string `form:"nickName" json:"nickName"`
	FavoriteArtist uint   `form:"favoriteArtist" json:"favoriteArtist"`
}

type ErrorObject struct {
	ErrorPoint		string	`json:"errorPoint"`
	ErrorMessage	string	`json:"errorMessage"`
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
	var	errorObj	ErrorObject
	//var errorPoint	string
	var httpCode	int
	var checker		bool

	err := c.ShouldBind(&inputData)
	if err != nil {
		errorObj.ErrorPoint = permissions.AnalyzeErrorMessage(err.Error())
		errorObj.ErrorMessage = "허용 되지 않는 값 입니다."
		c.JSON(400, gin.H {
			"err": errorObj,
		})
		return
	}
	errorObj.ErrorPoint, httpCode, checker = permissions.IsEmpty(&inputData)
	if checker {
		errorObj.ErrorMessage = "비어 있는 값 입니다."
		c.JSON(httpCode, gin.H {
			"err": errorObj,
		})
		return
	}
	inputDataToUser(&user, inputData)
	errorObj.ErrorPoint, httpCode, checker = permissions.IsExist(&user)
	if checker {
		errorObj.ErrorMessage = "이미 존재한 값 입니다."
		c.JSON(httpCode, gin.H {
			"err": errorObj,
		})
		return
	}
	password := hash(user.Password)
	user.Password = password
	tx := configs.DB.Create(&user)
	if tx.Error != nil {
		errorObj.ErrorPoint, httpCode = permissions.FindErrorPoint(tx.Error)
		errorObj.ErrorMessage = "잘못된 값 입니다."
		c.JSON(httpCode, gin.H {
			"err": errorObj,
		})
		return
	}
	c.JSON(httpCode, "signUp success")
}
