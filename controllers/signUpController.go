package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/DuckBap/Duckbap-backend/models"
	"net/http"
)

func SignUp(c *gin.Context) {
	user, err := models.SignUpUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}
	c.JSON(http.StatusOK, user)
}