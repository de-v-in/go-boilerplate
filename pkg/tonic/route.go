package tonic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"

	"github.com/gin-gonic/gin"
)

const (
	Get     = 1
	Post    = 2
	Put     = 3
	Delete  = 4
	Patch   = 5
	Options = 6
	Head    = 7
)

type RouteSchema struct {
	Summary     string
	Description string
	Querystring any
	Params      any
	Body        any
	Response    map[int]any `json:"response"`
}

type RouteDef struct {
	Method  int8
	Url     string
	Handler func(c *gin.Context)
	Tags    []string
	Schema  RouteSchema
}

var apiSpec = make(map[string]any)

func InitSwagger() {
	apiSpec["openapi"] = "3.0.0"
	apiSpec["info"] = map[string]any{
		"title":   "Go CRUD Example",
		"version": "1.0.0",
	}
	apiSpec["components"] = map[string]any{
		"schemas": make(map[string]any),
	}
	apiSpec["paths"] = make(map[string]any)
}

func CreateRoutes(router *gin.RouterGroup, routeDefs []RouteDef) {
	basePath := router.BasePath()
	apiSpecPaths, _ := apiSpec["paths"].(map[string]any)

	for _, routeDef := range routeDefs {
		if routeDef.Url == "/" {
			routeDef.Url = ""
		} else if len(routeDef.Url) > 0 && routeDef.Url[0] != '/' {
			routeDef.Url = fmt.Sprintf("/%s", routeDef.Url)
		}

		httpMethod := ""

		switch routeDef.Method {
		case Get:
			httpMethod = "get"
			router.GET(routeDef.Url, routeDef.Handler)
		case Post:
			httpMethod = "post"
			router.POST(routeDef.Url, routeDef.Handler)
		case Put:
			httpMethod = "put"
			router.PUT(routeDef.Url, routeDef.Handler)
		case Delete:
			httpMethod = "delete"
			router.DELETE(routeDef.Url, routeDef.Handler)
		case Patch:
			httpMethod = "patch"
			router.PATCH(routeDef.Url, routeDef.Handler)
		case Options:
			httpMethod = "options"
			router.OPTIONS(routeDef.Url, routeDef.Handler)
		case Head:
			httpMethod = "head"
			router.HEAD(routeDef.Url, routeDef.Handler)
		}

		apiPath := toSwaggerAPIPath(fmt.Sprintf("%s%s", basePath, routeDef.Url))

		pathSpecRaw, apiPathExisted := apiSpecPaths[apiPath]

		if !apiPathExisted {
			pathSpecRaw = make(map[string]any)
			apiSpecPaths[apiPath] = pathSpecRaw
		}

		pathSpec, _ := pathSpecRaw.(map[string]any)
		handlerSpec := make(map[string]any)
		if routeDef.Schema.Summary != "" {
			handlerSpec["summary"] = routeDef.Schema.Summary
		}
		if routeDef.Schema.Description != "" {
			handlerSpec["description"] = routeDef.Schema.Description
		}
		if routeDef.Tags != nil {
			handlerSpec["tags"] = routeDef.Tags
		}

		handlerSpec["parameters"] = []map[string]any{}
		if routeDef.Schema.Params != nil {
			// Param is always a struct
			paramType := reflect.TypeOf(routeDef.Schema.Params)
			paramObjSchema := ToSwaggerSchema(paramType)
			paramSchema := map[string]any{}
			for propName, propSchema := range paramObjSchema["properties"].(map[string]any) {
				paramSchema["in"] = "path"
				paramSchema["name"] = propName
				paramSchema["required"] = true
				paramSchema["schema"] = propSchema

				handlerSpec["parameters"] = append(handlerSpec["parameters"].([]map[string]any), paramSchema)
			}
		}
		if routeDef.Schema.Querystring != nil {
			// Querystring is always a struct
			queryType := reflect.TypeOf(routeDef.Schema.Querystring)
			querySchema := ToSwaggerSchema(queryType)
			for propName, propSchema := range querySchema["properties"].(map[string]any) {
				propSchema.(map[string]any)["in"] = "query"
				propSchema.(map[string]any)["name"] = propName
				propSchema.(map[string]any)["required"] = slices.Contains(querySchema["required"].([]string), propName)
				propSchema.(map[string]any)["schema"] = propSchema

				handlerSpec["parameters"] = append(handlerSpec["parameters"].([]any), propSchema)
			}
		}

		if routeDef.Schema.Body != nil {
			t := reflect.TypeOf(routeDef.Schema.Body)
			handlerSpec["requestBody"] = map[string]any{
				"content": map[string]any{
					"application/json": map[string]any{
						"schema": ToSwaggerSchema(t),
					},
				},
				"required": true,
			}
		}
		if routeDef.Schema.Response != nil {
			handlerSpec["responses"] = map[int]any{}
			for status, response := range routeDef.Schema.Response {
				respType := reflect.TypeOf(response)
				handlerSpec["responses"].(map[int]any)[status] = map[string]any{
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ToSwaggerSchema(respType),
						},
					},
					"description": "Default description",
				}
			}
		}

		pathSpec[httpMethod] = handlerSpec
	}
}

func GetApiSpecs() []byte {
	b, _ := json.Marshal(apiSpec)
	return b
}

func toSwaggerAPIPath(path string) string {
	modifiedPath := regexp.MustCompile(`/:(\w+)`).ReplaceAllStringFunc(path, func(match string) string {
		return fmt.Sprintf("/{%s}", match[2:])
	})

	return modifiedPath
}