// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package swagger

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
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/accounts": {
            "post": {
                "description": "Creates an bank account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Creates an account",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.CreateAccountRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/accounts.CreateAccountResponse"
                        }
                    },
                    "400": {
                        "description": "Could not parse request"
                    },
                    "409": {
                        "description": "Account already registered"
                    },
                    "422": {
                        "description": "Could not create account"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/accounts/{account_id}": {
            "get": {
                "description": "Retrieve an account by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Gets an account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID",
                        "name": "account_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/accounts.GetAccountResponse"
                        }
                    },
                    "400": {
                        "description": "Could not parse request"
                    },
                    "404": {
                        "description": "Account not found"
                    },
                    "422": {
                        "description": "Could not create account"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "accounts.CreateAccountRequest": {
            "type": "object",
            "required": [
                "document_number"
            ],
            "properties": {
                "credit_limit": {
                    "type": "integer",
                    "example": 15000
                },
                "document_number": {
                    "type": "string",
                    "example": "12345678900"
                }
            }
        },
        "accounts.CreateAccountResponse": {
            "type": "object",
            "properties": {
                "account_id": {
                    "type": "string"
                }
            }
        },
        "accounts.GetAccountResponse": {
            "type": "object",
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "available_credit_limit": {
                    "type": "integer"
                },
                "balance": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "document_number": {
                    "type": "string"
                },
                "updated_at": {
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
	Version:     "1.0",
	Host:        "localhost:3000",
	BasePath:    "/api/v1",
	Schemes:     []string{"http"},
	Title:       "Swagger Mybank API",
	Description: "Documentation Mybank API",
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
