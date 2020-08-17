// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/accounts/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "loginInfo",
                        "name": "loginInfo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/middlewares.Login"
                        }
                    }
                ],
                "responses": {
                    "200": {}
                }
            }
        },
        "/accounts/sign-up": {
            "post": {
                "description": "\u003cbr\u003e유저의 정보를 context에 저장하여 미들웨어에게 전달합니다.\u003cbr\u003e\n빈 값이 요청되었을 경우 오류 발생위치와 오류 메시지를 반환합니다.\u003cbr\u003e\n잘못된 값이 들어왔을 경우 오류 발생위치와 오류 메시지를 반환합니다.\u003cbr\u003e\n이미 회원인 경우 오류를 발생시켜 오류 발생위치와 오류 메시지를 반환합니다.\u003cbr\u003e",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "회원가입 요청이 들어왔을 때 동작하는 곳",
                "parameters": [
                    {
                        "description": "user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.InputUserData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "token",
                        "schema": {
                            "type": ""
                        },
                        "headers": {
                            "user": {
                                "type": "object",
                                "description": "token"
                            }
                        }
                    },
                    "208": {
                        "description": "이미 존재한 값이 들어올 때\", \"{\"err\": {\"errorPoint\": \"message\"}}",
                        "schema": {
                            "type": ""
                        }
                    },
                    "400": {
                        "description": "잘못된 값이 들어올 때\", \"{\"err\": {\"errorPoint\": \"message\"}}",
                        "schema": {
                            "type": ""
                        }
                    },
                    "404": {
                        "description": "해당 값을 통해서 회원 가입을 못할 때\", \"{\"err\": {\"errorPoint\": \"message\"}}",
                        "schema": {
                            "type": ""
                        }
                    },
                    "424": {
                        "description": "참조할 수 없는 값이 들어올 때\", \"{\"err\": {\"errorPoint\": \"message\"}}",
                        "schema": {
                            "type": ""
                        }
                    }
                }
            }
        },
        "/artists": {
            "get": {
                "description": "## \u003cbr\u003e아티스트 리스트를 반환합니다.\n\u003cbr\u003e\n# 쿼리스트링이 존재하지 않는 경우\n1. 모든 아티스트를 반환합니다.\u003cbr\u003e\n# \u003cbr\u003e쿼리스트링이 존재하는 경우\n1. 쿼리스트링을 조건으로 필터링 된 아티스트를 반환합니다.\u003cbr\u003e\n1. 회사에 속한 아티스트들 /v1/artists?ent-id=()",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "아티스트 리스트",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.OutputArtistList"
                            }
                        }
                    }
                }
            }
        },
        "/ents": {
            "get": {
                "description": "\u003cbr\u003e엔터테인먼트 리스트를 반환합니다.\u003cbr\u003e",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "엔터테인먼트 리스트",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.Entertainment"
                            }
                        }
                    }
                }
            }
        },
        "/fundings": {
            "get": {
                "description": "\u003cbr\u003e아티스트와 관련된 펀딩 리스트를 반환합니다.\n\u003cbr\u003e",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "아티스트와 관련 펀딩 리스트",
                "parameters": [
                    {
                        "type": "integer",
                        "description": ".",
                        "name": "artist-id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.fundingListResBody"
                            }
                        }
                    }
                }
            }
        },
        "/fundings/banner": {
            "get": {
                "description": "\u003cbr\u003e펀딩 리스트를 반환합니다.\u003cbr\u003e",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "메인 배너에서 보여줄 펀딩 리스트",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.bannerFunding"
                            }
                        }
                    }
                }
            }
        },
        "/fundings/main": {
            "get": {
                "description": "## 메인 페이지에서 보여줄 펀딩 리스트를 반환합니다.\n\u003cbr\u003e\n## 로그인이 되어있을 경우\n1. 즐겨찾기에 저장되어있는 아이돌, 최애 아이돌과 관련된 펀딩들, 판매량이 가장 높은 펀딩들이 포함됩니다.\n\u003cbr\u003e\n2. 펀딩 8개가 들어있는 리스트가 반환됩니다.\u003cbr\u003e\n\u003cbr\u003e\n## 로그인이 되어있지 않을 경우\n1. 판매량이 높은 펀딩 8개가 포함된 리스트가 반환됩니다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "메인에서 보여줄 펀딩 리스트",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.listFunding"
                            }
                        }
                    },
                    "404": {
                        "description": "해당 페이지에 대한 값을 찾을 수 없을 때, \"{\"err\": \"해당 값을 찾을 수 없습니다.\"}",
                        "schema": {
                            "type": ""
                        }
                    }
                }
            }
        },
        "/fundings/{fund_id}": {
            "get": {
                "description": "\u003cbr\u003e펀딩 상세정보를 반환합니다.\n\u003cbr\u003e\nsellerName : 판매자의 닉네임\u003cbr\u003e\nfundName : 펀드 이름\u003cbr\u003e\nprice : 하나를 구매할 때의 가격\u003cbr\u003e\ntargetAmount : 판매 목표량\u003cbr\u003e\nsalesAmount : 현재까지의 판매량\u003cbr\u003e\nstartDate : 펀딩 시작 일\u003cbr\u003e\nendDate: 펀딩 마감 일\u003cbr\u003e\nartistName : 펀딩과 관련된 연예인 이름\u003cbr\u003e\nachievementRate : 펀딩 달성률 (판매량 / 목표량)\u003cbr\u003e\ndDay : 펀딩 마감일까지 남은 날짜\u003cbr\u003e\nfundingImgUrls : 펀딩 상품의 이미지 주소들\u003cbr\u003e\ndetailedImgUrl : 펀딩 상세정보 이미지\u003cbr\u003e",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "펀딩 상세정보",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "fund_id",
                        "name": "fund_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.fundingResBody"
                        }
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "description": "로그인이 되어있어야 접근 가능합니다.\u003cbr\u003e\nnickName : \u003cbr\u003e\nfavoriteArtist : \u003cbr\u003e\nbuy : \u003cbr\u003e\nsell : \u003cbr\u003e\nbookmark : \u003cbr\u003e",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "마이 페이지에서 보여줄 정보",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.data"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Entertainment": {
            "type": "object",
            "properties": {
                "entId": {
                    "type": "integer"
                },
                "imgUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.InputUserData": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "favoriteArtist": {
                    "type": "integer"
                },
                "nickName": {
                    "type": "string"
                },
                "password1": {
                    "type": "string"
                },
                "password2": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        },
        "controllers.OutputArtistList": {
            "type": "object",
            "properties": {
                "artistId": {
                    "type": "integer"
                },
                "artistImgUrl": {
                    "type": "string"
                },
                "artistName": {
                    "type": "string"
                },
                "entId": {
                    "type": "integer"
                }
            }
        },
        "controllers.artist": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "imgUrl": {
                    "type": "string"
                },
                "level": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.bannerFunding": {
            "type": "object",
            "properties": {
                "artistID": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "mainImgUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.data": {
            "type": "object",
            "properties": {
                "bookmark": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.artist"
                    }
                },
                "buy": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.receipt"
                    }
                },
                "favoriteArtist": {
                    "$ref": "#/definitions/controllers.artist"
                },
                "id": {
                    "type": "integer"
                },
                "nickName": {
                    "type": "string"
                },
                "sell": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.receipt"
                    }
                }
            }
        },
        "controllers.funding": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "mainImgUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.fundingListResBody": {
            "type": "object",
            "properties": {
                "achievementRate": {
                    "type": "number"
                },
                "dDay": {
                    "type": "integer"
                },
                "fundingName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "mainImgUrl": {
                    "type": "string"
                },
                "sellerName": {
                    "type": "string"
                }
            }
        },
        "controllers.fundingResBody": {
            "type": "object",
            "properties": {
                "achievementRate": {
                    "description": "salesAmount / Price",
                    "type": "number"
                },
                "artistName": {
                    "type": "string"
                },
                "dDay": {
                    "type": "integer"
                },
                "detailedImgUrl": {
                    "type": "string"
                },
                "endDate": {
                    "type": "string"
                },
                "fundName": {
                    "type": "string"
                },
                "fundingImgUrls": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "salesAmount": {
                    "type": "integer"
                },
                "sellerName": {
                    "type": "string"
                },
                "startDate": {
                    "type": "string"
                },
                "targetAmount": {
                    "type": "integer"
                }
            }
        },
        "controllers.listFunding": {
            "type": "object",
            "properties": {
                "achievementRate": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "mainImgUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.receipt": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "funding": {
                    "$ref": "#/definitions/controllers.funding"
                },
                "fundingId": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "middlewares.Login": {
            "type": "object",
            "required": [
                "password",
                "userName"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}