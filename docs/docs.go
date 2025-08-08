package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "Crypto Price Service",
    "description": "Сервис для получения цены криптовалют, добавления и удаления монет.",
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "schemes": ["http"],
  "paths": {
    "/currency/price": {
      "post": {
        "summary": "Получить цену монеты",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "coin": { "type": "string", "example": "bitcoin" },
                "timestamp": { "type": "integer", "example": 1754603100 }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "examples": {
              "application/json": {
                "status": "OK",
                "data": {
                  "coin": "bitcoin",
                  "price": 117200,
                  "currency": "USD",
                  "timestamp": 1754603017
                }
              }
            }
          },
          "404": {
            "description": "Цена не найдена",
            "examples": {
              "application/json": { "status": "NotFound", "message": "Price not found" }
            }
          },
          "400": {
            "description": "Некорректный запрос",
            "examples": {
              "application/json": { "status": "ERROR", "message": "Invalid request" }
            }
          },
		  "500": {
            "description": "Ошибка сервиса при получении данных",
            "schema": { "type": "object", "example": { "status": "ERROR", "message": "Error while receiving data" } }
          }
        }
      }
    },
    "/currency/add": {
      "post": {
        "summary": "Добавить монету",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "name_coin": { "type": "string", "example": "bitcoin" }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Монета добавлена",
            "examples": {
              "application/json": { "status": "OK", "message": "Coin 'bitcoin' added", "data": {} }
            }
          },
          "400": {
            "description": "Ошибка",
            "examples": {
              "application/json": [
                { "status": "ERROR", "message": "Coin not found" },
                { "status": "ERROR", "message": "Incorrect input data" }
              ]
            }
          },
		  "500": {
            "description": "Ошибка сервиса при добавлении",
            "schema": { "type": "object", "example": { "status": "ERROR", "message": "Service error while adding" } }
          }
        }
      }
    },
    "/currency/remove": {
      "delete": {
        "summary": "Удалить монету",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "name_coin": { "type": "string", "example": "bitcoin" }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Монета удалена",
            "examples": {
              "application/json": { "status": "OK", "message": "Coin deleted", "data": {} }
            }
          },
          "400": {
            "description": "Ошибка",
            "examples": {
              "application/json": [
                { "status": "ERROR", "message": "This coin is not on the list" },
                { "status": "ERROR", "message": "Incorrect input data" }
              ]
            }
          },
		  "500": {
            "description": "Ошибка сервиса при удалении",
            "schema": { "type": "object", "example": { "status": "ERROR", "message": "Service error while deleting" } }
          }
        }
      }
    }
  }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
