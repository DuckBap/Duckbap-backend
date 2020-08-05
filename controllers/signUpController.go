package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/DuckBap/Duckbap-backend/models"
	"net/http"
)

func SignUp(c *gin.Context) {
	user, errorValue, define := models.SignUpUser(c)
	if define == 1 {
		c.JSON(http.StatusAlreadyReported, errorValue)
	} else if define == 2 {
		c.JSON(http.StatusFailedDependency, errorValue)
	} else if define == 3 {
		c.JSON(http.StatusNotFound, errorValue)
	} else {
		c.JSON(http.StatusOK, user)
	}
}