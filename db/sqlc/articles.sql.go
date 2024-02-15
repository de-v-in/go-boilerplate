// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: articles.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createArticle = `-- name: CreateArticle :one
INSERT INTO articles (title, content) VALUES ($1, $2) RETURNING "id"
`

type CreateArticleParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createArticle, arg.Title, arg.Content)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getArticleById = `-- name: GetArticleById :one
SELECT id, title, content, created_at, updated_at FROM articles WHERE id = $1
`

func (q *Queries) GetArticleById(ctx context.Context, id uuid.UUID) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticleById, id)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getArticles = `-- name: GetArticles :many
SELECT id, title, content, created_at, updated_at FROM articles
`

func (q *Queries) GetArticles(ctx context.Context) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, getArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Article{}
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}