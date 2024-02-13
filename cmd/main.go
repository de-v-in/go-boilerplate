package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	dbPkg "github.com/phucvinh57/go-crud-example/db"
	sqlc "github.com/phucvinh57/go-crud-example/db/sqlc"
	"github.com/phucvinh57/go-crud-example/internal/app/controllers"
	jsonschema "github.com/phucvinh57/go-crud-example/internal/pkg"
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
	article := router.Group("articles")
	{
		schema := jsonschema.GenJSONSchema(controllers.PostDTO{})
		fmt.Println(schema)

		ctrler := controllers.NewArticleCrtler(db, ctx)
		article.GET("", ctrler.GetArticles)
	}
}

func main() {
	initServer()
	setupRoutes()
	app.Run(":8080")
}
