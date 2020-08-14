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

func	changeString(str string, forQuery bool) string {
	var index	int
	var checker bool
	var newStr	string

	if forQuery {
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
	} else {
		rune := str[0]
		if rune >= 'A' && rune <= 'Z' {
			newStr = str[:1]
			newStr = strings.ToLower(newStr)
			newStr = newStr + str[1:]
		}
	}
	return newStr
}
//--------------------------------------------------------------------------------
func	AnalyzeErrorMessage(message string) string {
	var	errorPoint	string

	if strings.Contains(message, "users.user_name") {
		errorPoint = "userName"
	} else if strings.Contains(message, "users.email") {
		errorPoint = "email"
	} else if strings.Contains(message, "users.nick_name") {
		errorPoint = "nickName"
	} else if strings.Contains(message, "favorite_artist") {
		errorPoint = "favoriteArtist"
	} else if strings.Contains(message, "userName") {
		errorPoint = "userName"
	} else if strings.Contains(message, "email") {
		errorPoint = "email"
	} else if strings.Contains(message, "nickName") {
		errorPoint = "nickName"
	} else if strings.Contains(message, "favoriteArtist") {
		errorPoint = "favoriteArtist"
	} else {
		errorPoint = "error"
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

func	FindErrorPoint(err error) (map[string]string, int){
	var errorPoint	string
	var httpCode	int

	errorStruct, marshalingError := makeErrorStruct(err)
	errorMap := make(map[string]string)
	if marshalingError != nil {
		errorPoint = "json marshaling"
	} else {
		errorPoint = AnalyzeErrorMessage(errorStruct.Message)
	}
	if errorStruct.Number == 1062 {
		httpCode = http.StatusAlreadyReported
		errorMap[errorPoint] = "이미 존재한 값 입니다."
	} else if errorStruct.Number == 1452 {
		httpCode = http.StatusFailedDependency
		errorMap[errorPoint] = "참조할 수 없는 값 입니다."
	} else {
		httpCode = http.StatusNotFound
		errorMap[errorPoint] = "해당 값을 찾을 수 없습니다."
	}
	return errorMap, httpCode
}

/*
*** 유효한 값이 들어왔는지 확인하는 함수들
 */
func	IsAlreadyValuePresent(model interface{}, query string, value interface{}) bool {
	var	count			int64
	var	changeChecker	bool
	var queryString		string

	queryString = changeString(query, true)
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

func	isEmptyValue (elements reflect.Value, index int) (map[string]string, bool) {
	var	emptyPoint		string
	var isEmpty			bool

	errorMap := make(map[string]string)
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
		emptyPoint = changeString(emptyPoint, false)
		errorMap[emptyPoint] = "비어 있는 값 입니다."
	}
	return errorMap, isEmpty
}

func	isAlreadyPresent(dataStruct interface{}, elements reflect.Value, index int) (map[string]string, bool){
	var isExist			bool
	var	model			interface{}

	errorMap := make(map[string]string)
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
	presentPoint = changeString(presentPoint, false)
	if isExist {
		errorMap[presentPoint] = "이미 존재한 값 입니다."
	}
	return errorMap, isExist
}

func	isImpossibleValue(elements reflect.Value, index int) (map[string]string, bool) {
	var	value			string
	var	nextValue		string
	var impossiblePoint string
	var	isImpossible	bool
	var idx				int

	idx = index
	elementName := elements.Type().Field(idx).Name
	errorMap := make(map[string]string)
	if elementName == "Password2" {
		if elements.Type().Field(idx - 1).Name == "Password1" {
			value = fmt.Sprintf("%v",elements.Field(idx).Interface())
			nextValue = fmt.Sprintf("%v",elements.Field(idx - 1).Interface())
			if value != nextValue {
				impossiblePoint = elementName
				errorMap[impossiblePoint] = "비밀번호가 일치하지 않습니다."
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
			errorMap[impossiblePoint] = "잘못된 형식 입니다."
		}
	}
	return errorMap, isImpossible
}

func	IsEmpty (dataStruct interface{}) (map[string]string, int, bool) {
	var isPossible		bool
	var	httpCode		int
	var	errorChecker	bool

	errorMap := make(map[string]string)
	target := reflect.ValueOf(dataStruct)
	elements := target.Elem()
	for idx := 0; idx < elements.NumField(); idx++ {
		 if emptyValue,isExist := isEmptyValue(elements, idx); isExist {
			errorMap = emptyValue
			isPossible = isExist
			errorChecker = true
			break
		} else if impossibleValue, isUnable := isImpossibleValue(elements, idx); isUnable {
			errorMap = impossibleValue
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
	return errorMap, httpCode, isPossible
}

func	IsExist(dataStruct interface{}) (map[string]string, int, bool){
	var isPossible		bool
	var	httpCode		int

	errorMap := make(map[string]string)
	target := reflect.ValueOf(dataStruct)
	elements := target.Elem()
	httpCode = http.StatusOK
	for idx := 0; idx < elements.NumField(); idx++ {
		presentPoint, isExist := isAlreadyPresent(dataStruct, elements, idx)
		if isExist {
			errorMap = presentPoint
			isPossible = isExist
			httpCode = http.StatusAlreadyReported
			break
		}
	}
	return errorMap, httpCode, isPossible
}
