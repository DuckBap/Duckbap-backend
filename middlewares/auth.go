package middlewares

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Auth *jwt.GinJWTMiddleware

type Login struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// @Summary 로그인.
// @Description <br>로그인 성공 시 jwt 토큰을 반환합니다.
// @Description <br>

// @Param loginInfo body Login true "loginInfo"
// @Accept  json
// @Produce  json
// @Router /accounts/login [post]
// @Success 200
func init() {
	Auth, _ = jwt.New(&jwt.GinJWTMiddleware{
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
			var login Login
			var user models.User
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := login.UserName
			password := login.Password
			err := configs.DB.Where("user_name = ?", username).First(&user).Error
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &user, nil
		},
	})
}
