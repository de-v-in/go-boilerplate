package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	dbPkg "github.com/phucvinh57/go-crud-example/db"
	sqlc "github.com/phucvinh57/go-crud-example/db/sqlc"
	"github.com/phucvinh57/go-crud-example/internal/app/controllers"
	"github.com/phucvinh57/go-crud-example/pkg/tonic"
	"github.com/rs/zerolog"
)

var (
	app    *gin.Engine
	ctx    context.Context
	db     *sqlc.Queries
	router *gin.RouterGroup
)

func initServer() {
	ctx = context.Background()
	app = gin.Default()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	router = app.Group("/api")
	db = sqlc.New(dbPkg.Init())
}

func setupRoutes() {
	tonic.Init(&tonic.Config{
		OpenAPIVersion: "3.0.0",
		Info: map[string]interface{}{
			"title":       "Go CRUD Example",
			"description": "A simple CRUD example using Go and PostgreSQL",
			"version":     "1.0.0",
		},
	})

	article := router.Group("articles")
	{
		ctrler := controllers.NewArticleCrtler(db, ctx)
		routes := []tonic.Route{
			{
				Method: tonic.Get,
				Url:    "",
				Schema: &tonic.RouteSchema{
					Summary: "Get all articles",
					Response: map[int]interface{}{
						200: []controllers.ArticleDTO{},
					},
				},
				Handler: ctrler.GetArticles,
			},
			{
				Method: tonic.Post,
				Url:    "",
				Schema: &tonic.RouteSchema{
					Summary: "Create an article",
					Body:    controllers.ArticleMutationDTO{},
					Response: map[int]interface{}{
						200: struct {
							ID string `json:"id"`
						}{},
					},
				},
				Handler: ctrler.CreateArticle,
			},
			{
				Method: tonic.Get,
				Url:    ":id",
				Schema: &tonic.RouteSchema{
					Params: struct {
						ID string `json:"id" binding:"required"`
					}{},
					Response: map[int]interface{}{
						200: controllers.ArticleDTO{},
					},
				},
				Handler: ctrler.GetArticleById,
			},
		}
		for i := range routes {
			routes[i].Tags = []string{"articles"}
		}
		tonic.CreateRoutes(article, routes)
	}
}

func main() {
	initServer()

	tonic.Init(&tonic.Config{
		OpenAPIVersion: "3.0.0",
		Info: map[string]interface{}{
			"title":       "Go CRUD Example",
			"description": "A simple CRUD example using Go and PostgreSQL",
			"version":     "1.0.0",
		},
	})

	setupRoutes()

	app.GET("/docs/*w", gin.WrapH(http.StripPrefix("/docs", tonic.GetHandler())))
	app.Run(":8080")
}
