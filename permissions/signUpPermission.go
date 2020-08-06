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

/*
*** 유효한 값이 들어왔는지 확인하는 함수들
 */
func	IsAlreadyValuePresent(model interface{}, query string, value interface{}) bool {
	var	count			int64
	var	changeChecker	bool
	var queryString		string

	queryString = changeString(query)
	if queryString == "favorite_artist" {
		queryString = "id = ?"
		changeChecker = true
	} else {
		queryString += " = ?"
	}
	configs.DB.Model(model).Where(queryString, value).Count(&count)
	if !changeChecker && count != 0 {
		return true
	} else if changeChecker && count == 0 {
		return true
	}
	return false
}

func	IsPossibleValue(dataStruct interface{}) (string, int, bool){
	var	emptyPoint		string
	var elementTag		string
	var emptyBool		bool
	var elementString	string
	var elementLen		int
	var	httpCode		int
	var	model			interface{}
	var elementField	interface{}

	target := reflect.ValueOf(dataStruct)
	elements := target.Elem()
	httpCode = http.StatusOK
	for idx := 0; idx < elements.NumField(); idx++ {
		elementField = elements.Field(idx).Interface()
		fmt.Printf("elementField : %v\n", elements.Field(idx))
		fmt.Printf("another elementField : %v\n", elements.Type())
		fmt.Printf("another elementField2 : %v\n", elements.Type().Field(idx))
		elementString = fmt.Sprintf("%v", elementField)
		elementLen = len(elementString)
		elementTag = elements.Type().Field(idx).Tag.Get("gorm")
		emptyPoint = elements.Type().Field(idx).Name
		if  strings.Contains(elementTag, "not null") {
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
				if strings.Contains(elementTag, "unique") &&
					IsAlreadyValuePresent(model, emptyPoint, elementField) {
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