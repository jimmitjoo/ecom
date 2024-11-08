// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/products": {
            "get": {
                "description": "Retrieves a list of all products, sorted by creation date in descending order (newest first)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "List all products",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number (starts from 1)",
                        "name": "page",
                        "in": "query",
                        "default": 1
                    },
                    {
                        "type": "integer",
                        "description": "Number of products per page",
                        "name": "size",
                        "in": "query",
                        "default": 10
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {
                                        "$ref": "#/definitions/models.Product"
                                    }
                                },
                                "pagination": {
                                    "type": "object",
                                    "properties": {
                                        "current_page": {
                                            "type": "integer",
                                            "example": 1
                                        },
                                        "page_size": {
                                            "type": "integer",
                                            "example": 10
                                        },
                                        "total_items": {
                                            "type": "integer",
                                            "example": 100
                                        },
                                        "total_pages": {
                                            "type": "integer",
                                            "example": 10
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new product with the specified details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Create a new product",
                "parameters": [
                    {
                        "description": "Product details",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/products/batch": {
            "put": {
                "description": "Updates multiple products in a single request. All products must exist and contain valid data.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Batch update multiple products simultaneously",
                "parameters": [
                    {
                        "description": "Array of products to update with their IDs and new data",
                        "name": "products",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Product"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Array of updated products",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Product"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid JSON data or validation errors",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "404": {
                        "description": "One or more products not found",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates multiple products simultaneously in a single request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Create multiple products in bulk",
                "parameters": [
                    {
                        "description": "Array of products to create",
                        "name": "products",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Product"
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Array of created products",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Product"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid JSON data",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes multiple products in a single request by their IDs. Returns results of deletion operations.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Batch delete multiple products simultaneously",
                "parameters": [
                    {
                        "description": "Array of product IDs to delete",
                        "name": "productIDs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Map of product IDs to deletion status",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid JSON data",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "404": {
                        "description": "One or more products not found",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.APIError"
                        }
                    }
                }
            }
        },
        "/products/{id}": {
            "get": {
                "description": "Retrieves a product with the specified ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates an existing product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Update a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated product details",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a product with the specified ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Delete a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "Invalid request"
                }
            }
        },
        "models.APIError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.MarketMetadata": {
            "type": "object",
            "required": [
                "market",
                "title"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "keywords": {
                    "type": "string"
                },
                "market": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Price": {
            "type": "object",
            "required": [
                "amount",
                "currency"
            ],
            "properties": {
                "amount": {
                    "type": "number",
                    "minimum": 0
                },
                "currency": {
                    "type": "string"
                }
            }
        },
        "models.Product": {
            "type": "object",
            "required": [
                "base_title",
                "id",
                "metadata",
                "prices",
                "sku"
            ],
            "properties": {
                "base_title": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_hash": {
                    "description": "Hash of last known state",
                    "type": "string"
                },
                "metadata": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.MarketMetadata"
                    }
                },
                "prices": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Price"
                    }
                },
                "sku": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "variants": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Variant"
                    }
                },
                "version": {
                    "description": "Version number for optimistic locking",
                    "type": "integer"
                }
            }
        },
        "models.Stock": {
            "type": "object",
            "required": [
                "location_id"
            ],
            "properties": {
                "location_id": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer",
                    "minimum": 0
                }
            }
        },
        "models.Variant": {
            "type": "object",
            "required": [
                "attributes",
                "id",
                "sku"
            ],
            "properties": {
                "attributes": {
                    "description": "e.g. {\"size\": \"XL\", \"color\": \"blue\"}",
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "sku": {
                    "type": "string"
                },
                "stock": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Stock"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "ws"},
	Title:            "E-commerce Product API",
	Description:      "A robust and scalable API for product management in e-commerce systems",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
