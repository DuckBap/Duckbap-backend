package permissions

import (
	"encoding/json"
	"fmt"
	"github.com/DuckBap/Duckbap-backend/configs"
	"net/http"
	"reflect"
	"strings"
)

type ErrorStruct struct {
	Number	uint
	Message	string
}

func	changeString(str string) string {
	var index	int
	var checker bool
	var newStr	string

	for idx,rune := range str {
		if idx != 0 && (rune >= 'A' && rune <= 'Z') {
			index = idx
			checker = true
			break
		}
	}
	if checker {
		newStr = str[:index] + "_" + str[index:]
	} else {
		newStr = str
	}
	newStr = strings.ToLower(newStr)
	return newStr
}

func	analyzeErrorMessage(message string) string {
	var	errorPoint	string

	if strings.Contains(message, "users.user_name") {
		errorPoint = "user_name"
	} else if strings.Contains(message, "users.email") {
		errorPoint = "email"
	} else if strings.Contains(message, "users.nick_name") {
		errorPoint = "nick_name"
	} else if strings.Contains(message, "favorite_artist") {
		errorPoint = "favorite_artist"
	}
	return errorPoint
}

func	makeErrorStruct(err error) (ErrorStruct, error) {
	var	errorStruct	ErrorStruct
	var	marshalingError	error

	errorJson,marshalError := json.Marshal(&err)
	unmarshalError := json.Unmarshal(errorJson, &errorStruct)
	if marshalError != nil {
		marshalingError = marshalError
	} else if unmarshalError != nil {
		marshalingError = unmarshalError
	}
	return errorStruct, marshalingError
}

func	FindErrorPoint(err error) (string, int){
	var errorPoint 	string
	var httpCode	int

	errorStruct, marshalingError := makeErrorStruct(err)
	if marshalingError != nil {
		errorPoint = "json marshaling error"
	} else {
		errorPoint = analyzeErrorMessage(errorStruct.Message)
	}
	if errorStruct.Number == 1062 {
		httpCode = http.StatusAlreadyReported
	} else if errorStruct.Number == 1452 {
		httpCode = http.StatusFailedDependency
	} else {
		httpCode = http.StatusNotFound
	}
	return errorPoint, httpCode
}

func	IsAlreadyValuePresent(model interface{}, query string, value interface{}) bool {
	var	count			int64
	var	changeChecker	bool

	query = changeString(query)
	if query == "favorite_artist" {
		configs.DB.Model(model).Where("id = ?", value).Count(&count)
		changeChecker = true
	} else {
		query += " = ?"
		configs.DB.Model(model).Where(query, value).Count(&count)
	}
	if !changeChecker && count != 0 {
		return true
	} else if changeChecker && count == 0 {
		return true
	}
	return false
}

func	IsEmptyValue(dataStruct interface{}) (string, int, bool){
	var	emptyPoint		string
	var emptyBool		bool
	var elementString	string
	var elementLen		int
	var	httpCode		int
	var	model			interface{}

	target := reflect.ValueOf(dataStruct)
	elements := target.Elem()
	httpCode = http.StatusOK
	for idx := 0; idx < elements.NumField(); idx++ {
		elementString = fmt.Sprintf("%v",elements.Field(idx).Interface())
		elementLen = len(elementString)
		elementType := elements.Type().Field(idx).Type.String()
		emptyPoint = elements.Type().Field(idx).Name
		if elementType == "uint" || elementType == "string" {
			if elementLen == 0 || elementString == "0" {
				httpCode = http.StatusBadRequest
				emptyBool = true
				break
			} else {
				if emptyPoint == "FavoriteArtist" {
					model = elements.Field(elements.NumField() - 1).Interface()
				} else {
					model = dataStruct
				}
				if IsAlreadyValuePresent(model, emptyPoint, elements.Field(idx).Interface()) {
					httpCode = http.StatusAlreadyReported
					emptyBool = true
					break
				}
			}
		}
	}
	emptyPoint = changeString(emptyPoint)
	return emptyPoint, httpCode, emptyBool
}