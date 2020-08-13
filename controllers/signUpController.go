package controllers

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/permissions"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type InputUserData struct {
	UserName       string `form:"userName" json:"userName"`
	//Password1      string `form:"password1" json:"password1"`
	Password      string `form:"password" json:"password"`
	Password2      string `form:"password2" json:"password2"`
	Email          string `form:"email" json:"email"`
	NickName       string `form:"nickName" json:"nickName"`
	FavoriteArtist uint   `form:"favoriteArtist" json:"favoriteArtist"`
}

func inputDataToUser(user *models.User, inputData InputUserData) {
	(*user).UserName = inputData.UserName
	(*user).Password = inputData.Password
	//(*user).Password = inputData.Password1
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
		return
	}
	errorMap, httpCode, checker = permissions.IsEmpty(&inputData)
	if checker {
		c.JSON(httpCode, gin.H {
			"err": errorMap,
		})
		return
	}
	inputDataToUser(&user, inputData)
	errorMap, httpCode, checker = permissions.IsExist(&user)
	if checker {
		c.JSON(httpCode, gin.H {
			"err": errorMap,
		})
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
		return
	}
	//token := autoLogin(user)
	//fmt.Println("string token  ",token)
	//c.SetCookie("user", token, 3600, "/", "localhost", false, true)
	//token, err := autoLogin(c, user)
	//if err != nil {
	//	return
	//}
	//c.JSON(httpCode, token)
	//c.JSON(httpCode, user)
	//b := a(c, user)
	//b.LoginHandler(c)
}

func a(c *gin.Context, user models.User) *jwt.GinJWTMiddleware{
	Auth, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Duckbap",
		Key:         []byte("NEED SECRET KEY"),
		IdentityKey: "user",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id": v.ID,
					"userName": v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Model: gorm.Model{ID:uint(claims["id"].(float64))},
				UserName: claims["userName"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return &user, nil
		},
	})
	return Auth
}
//func	autoLogin(userStruct models.User) string {
//	mySigningKey := []byte("NEED SECRET KEY")
//	tmp := jwt.MapClaims{
//		"id": userStruct.ID,
//		"userName": userStruct.UserName,
//	}
//	type MyCustomClaims struct {
//		Foo string `json:"foo"`
//		jwt.MapClaims
//	}
//	// Create the Claims
//	claims := MyCustomClaims{
//		"bar",
//		tmp,
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
//	fmt.Println("token :: ",token)
//	ss,_ := token.SignedString(mySigningKey)
//	fmt.Printf("string   %v",ss)
//	return ss
//}
