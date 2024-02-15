package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sqlc "github.com/phucvinh57/go-crud-example/db/sqlc"
)

type ArticleDTO struct {
	ID      string `json:"id"`
	Title   string `json:"title" binding:"required,min=4,max=255"`
	Content string `json:"content" binding:"required"`
}

type ArticleMutationDTO struct {
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

func (ac *ArticleCtrler) GetArticles(c *gin.Context) {
	posts, err := ac.db.GetArticles(ac.ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, posts)
}

func (ac *ArticleCtrler) CreateArticle(c *gin.Context) {
	var article ArticleMutationDTO
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	articleId, err := ac.db.CreateArticle(ac.ctx, sqlc.CreateArticleParams{
		Title:   article.Title,
		Content: article.Content,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": articleId})
}

func (ac *ArticleCtrler) GetArticleById(c *gin.Context) {
	articleIdStr := c.Param("id")
	articleId := uuid.MustParse(articleIdStr)
	article, err := ac.db.GetArticleById(ac.ctx, articleId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	articleDTO := ArticleDTO{
		ID:      article.ID.String(),
		Title:   article.Title,
		Content: article.Content,
	}
	c.JSON(200, articleDTO)
}
