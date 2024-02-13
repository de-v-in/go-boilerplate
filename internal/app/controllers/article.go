package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
	sqlc "github.com/phucvinh57/go-crud-example/db/sqlc"
)

type PostDTO struct {
	ID      int64  `json:"id"`
	Title   string `json:"title" binding:"required,min=4,max=255"`
	Content string `json:"content" binding:"required"`
}

type ArticleCtrler struct {
	db  *sqlc.Queries
	ctx context.Context
}

func NewArticleCrtler(db *sqlc.Queries, ctx context.Context) *ArticleCtrler {
	return &ArticleCtrler{db, ctx}
}

func (pc *ArticleCtrler) GetArticles(c *gin.Context) {
	posts, err := pc.db.GetArticles(pc.ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, posts)
}
