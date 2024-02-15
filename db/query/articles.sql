-- name: GetArticles :many
SELECT * FROM articles;

-- name: CreateArticle :one
INSERT INTO articles (title, content) VALUES ($1, $2) RETURNING "id";

-- name: GetArticleById :one
SELECT * FROM articles WHERE id = $1;
