package models

import (
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/permissions"
	"github.com/gin-gonic/gin"
)


func SignUpUser(c *gin.Context) (User, string, uint) {
	var	user	User
	var	errorMessage	string
	var define			uint
	var boolean			bool

	c.Bind(&user)
	errorMessage, define, boolean = permissions.CheckInputValue(user)
	if boolean {
		return user, errorMessage, define
	}
	errorMessage, define, boolean = permissions.FindValueAlreadyPresent(user)
	if boolean {
		return user, errorMessage, define
	}
	//
	//if err != nil {
	//	return user, "fail_binding", 3
	//}
	tx := configs.DB.Create(&user)
	errorMessage, define = permissions.MakeErrorMessage(tx)
	return user, errorMessage, define
}