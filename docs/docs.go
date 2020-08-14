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
                "responses": {
                    "200": {
                        "description": "token",
                        "schema": {
                            "type": ""
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
        "/artists/": {
            "get": {
                "description": "\u003cbr\u003e아티스트 리스트를 반환합니다.\u003cbr\u003e\n쿼리스트링이 존재하지 않을 경우 모든 아티스트를 반환합니다.\u003cbr\u003e\n쿼리스트링이 존재하는 경우 쿼리스트링을 조건으로 필터링 된 아티스트를 반환합니다.\u003cbr\u003e\n쿼리스트링 종류\n1. /v1/artists?ent-id=()",
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
        "/ents/": {
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
