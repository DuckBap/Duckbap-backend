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

func	AnalyzeErrorMessage(message string) string {
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
		errorPoint = AnalyzeErrorMessage(errorStruct.Message)
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

func	isEmptyValue (elements reflect.Value, index int) (string, bool) {
	var	emptyPoint		string
	var isEmpty			bool

	elementField := elements.Field(index).Interface()
	elementType := elements.Field(index).Type().String()
	elementString := fmt.Sprintf("%v", elementField)
	emptyPoint = elements.Type().Field(index).Name
	if  elementType == "uint" || elementType == "string"{
		if len(elementString) == 0 || elementString == "0" {
			isEmpty = true
		} else if len(strings.Trim(elementString, " ")) == 0 {
			isEmpty = true
		}
	}
	if isEmpty {
		emptyPoint = changeString(emptyPoint)
	}
	return emptyPoint, isEmpty
}

func	isAlreadyPresent(dataStruct interface{}, elements reflect.Value, index int) (string, bool){
	var isExist			bool
	var	model			interface{}

	elementField := elements.Field(index).Interface()
	elementTag := elements.Type().Field(index).Tag.Get("gorm")
	presentPoint := elements.Type().Field(index).Name
	if  strings.Contains(elementTag, "not null") {
		if presentPoint == "FavoriteArtist" {
			model = elements.Field(elements.NumField() - 1).Interface()
		} else {
			model = dataStruct
		}
		if strings.Contains(elementTag, "unique") &&
			IsAlreadyValuePresent(model, presentPoint, elementField) {
			isExist = true
		}
	}
	presentPoint = changeString(presentPoint)
	return presentPoint, isExist
}

func	isImpossibleValue(elements reflect.Value, index *int) (string, bool) {
	var	value			string
	var	nextValue		string
	var impossiblePoint string
	var	isImpossible	bool
	var idx				int

	idx = *index
	elementName := elements.Type().Field(idx).Name
	if elementName == "Password2" {
		if elements.Type().Field(idx - 1).Name == "Password1" {
			value = fmt.Sprintf("%v",elements.Field(idx).Interface())
			nextValue = fmt.Sprintf("%v",elements.Field(idx - 1).Interface())
			if value != nextValue {
				impossiblePoint = "different password"
				*index++
				isImpossible = true
			}
		}
	} else if elementName == "Email" {
		value = fmt.Sprintf("%v", elements.Field(idx).Interface())
		_index := strings.Index(value, "@")
		if strings.Count(value, "@") != 1 {
			isImpossible = true
		} else if strings.Index(value, ".com") == _index + 1 {
			isImpossible = true
		} else if strings.Index(value, ".net") == _index + 1 {
			isImpossible = true
		} else if strings.Count(value, ".com") + strings.Count(value, ".net") != 1 {
			isImpossible = true
		} else if strings.Index(value, ".com") != (len(value) - 4) {
			if strings.Index(value, ".net") != (len(value) - 4) {
				isImpossible = true
			}
		}
		if isImpossible {
			impossiblePoint = elementName
		}
	}
	return impossiblePoint, isImpossible
}

func	IsEmpty (dataStruct interface{}) (string, int, bool) {
	var	errorPoint		string
	var isPossible		bool
	var	httpCode		int
	var	errorChecker	bool

	target := reflect.ValueOf(dataStruct)
	elements := target.Elem()
	for idx := 0; idx < elements.NumField(); idx++ {
		 if emptyValue,isExist := isEmptyValue(elements, idx); isExist {
			errorPoint = emptyValue
			isPossible = isExist
			errorChecker = true
			break
		} else if impossibleValue, isUnable := isImpossibleValue(elements, &idx); isUnable {
			errorPoint = impossibleValue
			isPossible = isUnable
			errorChecker = true
			break
		}
	}
	if errorChecker {
		httpCode = http.StatusBadRequest
	} else {
		httpCode = http.StatusOK
	}
	return errorPoint, httpCode, isPossible
}

func	IsExist(dataStruct interface{}) (string, int, bool){
	var	errorPoint		string
	var isPossible		bool
	var	httpCode		int

	target := reflect.ValueOf(dataStruct)
	elements := target.Elem()
	httpCode = http.StatusOK
	for idx := 0; idx < elements.NumField(); idx++ {
		presentPoint, isExist := isAlreadyPresent(dataStruct, elements, idx)
		if isExist {
			errorPoint = presentPoint
			isPossible = isExist
			httpCode = http.StatusAlreadyReported
			break
		}
	}
	return errorPoint, httpCode, isPossible
}
