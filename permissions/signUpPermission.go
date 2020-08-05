package permissions

import (
	"encoding/json"
	"github.com/DuckBap/Duckbap-backend/configs"
	"github.com/DuckBap/Duckbap-backend/models"
	"gorm.io/gorm"
	"strings"
)

type ErrorStruct struct {
	Number	uint
	Message	string
}

func	findContainPart(message string) (string, uint) {
	var	returnString	string
	var	defineValue		uint

	if strings.Contains(message, "users.user_name") {
		returnString = "user_name"
		defineValue = 1
	} else if strings.Contains(message, "users.email") {
		returnString = "email"
		defineValue = 1
	} else if strings.Contains(message, "users.nick_name") {
		returnString = "nick_name"
		defineValue = 1
	} else if strings.Contains(message, "favorite_artist") {
		returnString = "favorite_artist"
		defineValue = 2
	}
	return returnString, defineValue
}

func	findErrorStruct(err error) (ErrorStruct, error) {
	var	errorStruct	ErrorStruct
	var	newError	error

	errorJson,marshalingError := json.Marshal(&err)
	unmarshalError := json.Unmarshal(errorJson, &errorStruct)
	if marshalingError != nil {
		newError = marshalingError
	} else if unmarshalError != nil {
		newError = unmarshalError
	}
	return errorStruct, newError
}

func	MakeErrorMessage(tx *gorm.DB) (string, uint){
	var errorPart 	string
	var define		uint

	if tx.Error != nil {
		errorStruct, err := findErrorStruct(tx.Error)
		if err != nil {
			errorPart = "json marshaling error"
			define = 3
			return errorPart, define
		}
		errorPart,define = findContainPart(errorStruct.Message)
		return errorPart, define
	}
	return errorPart, define
}

func	checkAlreadyValue(query string, value interface{}, model interface{}) bool {
	var	count	int64

	if query == "favorite_artist = ?" {
		configs.DB.Model(&models.Artist{}).Where("id = ?", value).Count(&count)

	} else {
		configs.DB.Model(model).Where(query, value).Count(&count)
	}
	if count != 0 {
		return true
	}
	return false
}

func	FindValueAlreadyPresent(user models.User) (string, uint, bool) {
	var	alreadyPresentValue string

	if checkAlreadyValue("user_name = ?", user.UserName, &user) {
		alreadyPresentValue = "user_name"
	} else if checkAlreadyValue("email = ?", user.Email, &user) {
		alreadyPresentValue = "email"
	} else if checkAlreadyValue("nick_name = ?", user.NickName, &user) {
		alreadyPresentValue = "nick_name"
	} else if checkAlreadyValue("favorite_artist = ?", user.FavoriteArtist, &user){
		alreadyPresentValue = "favorite_artist"
	} else {
		return alreadyPresentValue, 0, false
	}
	return alreadyPresentValue, 1, true
}

func	CheckInputValue(user models.User) (string, uint, bool) {
	var	errorValue	string
	if user.UserName == "" {
		errorValue = "user_name"
	} else if user.Password == "" {
		errorValue = "password"
	} else if user.NickName == "" {
		errorValue = "nick_name"
	} else if user.Email == "" {
		errorValue = "email"
	} else if user.FavoriteArtist == 0 {
		errorValue = "favorite_artist"
		return errorValue, 2, true
	} else {
		return errorValue, 0, false
	}
	return errorValue, 1, true
}