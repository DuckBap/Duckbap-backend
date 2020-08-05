package models

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/gin-gonic/gin"
)

func SignUpUser(c *gin.Context) (User, error) {
	var	user	User

	c.Bind(&user)
	tx := configs.DB.Create(&user)
	return user, tx.Error
}