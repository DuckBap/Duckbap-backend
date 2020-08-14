package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/permissions"
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

// @Summary 회원가입 요청이 들어왔을 때 동작하는 곳
// @Description <br>유저의 정보를 context에 저장하여 미들웨어에게 전달합니다.<br>
// @Description 빈 값이 요청되었을 경우 오류 발생위치와 오류 메시지를 반환합니다.<br>
// @Description 잘못된 값이 들어왔을 경우 오류 발생위치와 오류 메시지를 반환합니다.<br>
// @Description 이미 회원인 경우 오류를 발생시켜 오류 발생위치와 오류 메시지를 반환합니다.<br>
// @Accept  json
// @Produce  json
// @Router /accounts/sign-up [post]
// @Success 200 {} string "token"
// @Failure	208 {} string ""이미 존재한 값이 들어올 때", "{"err": {"errorPoint": "message"}}"
// @Failure	400 {} string ""잘못된 값이 들어올 때", "{"err": {"errorPoint": "message"}}"
// @Failure	404 {} string ""해당 값을 통해서 회원 가입을 못할 때", "{"err": {"errorPoint": "message"}}"
// @Failure	424 {} string ""참조할 수 없는 값이 들어올 때", "{"err": {"errorPoint": "message"}}"
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

func filterStruct(data *models.User) {
	(*data).FavoriteArtist = 0
	(*data).Password = ""
	(*data).NickName = ""
	(*data).Email = ""
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
	filterStruct(&user)
	c.Set("user", &user)
}
