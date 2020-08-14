package middlewares

import (
	"github.com/DuckBap/Duckbap-backend/models"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var	AutoLogin	*jwt.GinJWTMiddleware

func	init() {
	AutoLogin, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Duckbap",
		Key:         []byte("NEED SECRET KEY"),
		IdentityKey: "user",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":       v.ID,
					"userName": v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Model: gorm.Model{ID: uint(claims["id"].(float64))},
				UserName: claims["userName"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			userS, err := c.Get("user")
			if !err {
				return nil, jwt.ErrFailedAuthentication
			}
			user := userS.(*models.User)
			return user, nil
		},
	})
}