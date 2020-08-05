package middlewares

import (
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/appleboy/gin-jwt/v2"
)

var Auth *jwt.GinJWTMiddleware

func init() {
	Auth, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Duckbap",
		Key:         []byte("NEED SECRET KEY"),
		IdentityKey: "user",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id": v.ID,
					"username": v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Model: gorm.Model{ID:uint(claims["id"].(float64))},
				UserName: claims["username"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			type Login struct {
				UserName string `form:"username" json:"username" binding:"required"`
				Password string `form:"password" json:"password" binding:"required"`
			}
			var login Login
			var user models.User
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := login.UserName
			password := login.Password
			err := configs.DB.Where("username = ? AND password = ?", username, password).First(&user).Error
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &user, nil
		},
	})
}