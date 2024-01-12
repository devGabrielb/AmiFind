package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/devGabrielb/AmiFind/internal/entities"
)

type PostPaged struct {
	Page        int             `json:"page"`
	Total_pages int             `json:"total_pages"`
	Data        []entities.Post `json:"data"`
}
type PostRepository interface {
	GetPostPaged(ctx context.Context, page int) (*PostPaged, error)
}
type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}
func (pr *postRepository) Store(ctx context.Context, post entities.Post) (int64, error) {
	query, err := pr.db.PrepareContext(ctx, "INSERT INTO posts(title,description, date, pet_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	p, err := query.ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	postId, err := p.LastInsertId()
	if err != nil {
		return 0, err
	}
	return postId, nil
}

func (pr *postRepository) GetPostPaged(ctx context.Context, page int) (*PostPaged, error) {

	posts := make([]entities.Post, 0)

	limit := 10
	offset := limit * (page - 1)

	total, err := pr.totalOfpages(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	query := "SELECT id,title,description,date,pet_id FROM posts ORDER BY id LIMIT ? OFFSET ?"
	rows, err := pr.db.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		post := entities.Post{}

		if err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Date, &post.PetId); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	reponse := &PostPaged{Page: page, Total_pages: total, Data: posts}
	return reponse, nil
}

func (pr *postRepository) totalOfpages(ctx context.Context, page int, limit int) (int, error) {
	qtdPosts := 0

	err := pr.db.QueryRowContext(ctx, "SELECT count(id) FROM posts").Scan(&qtdPosts)
	if err != nil {
		return 0, err
	}
	log.Println("qtdPosts: ", qtdPosts)
	if qtdPosts <= limit {
		return 1, nil
	}
	total := (qtdPosts / limit)
	log.Println(total)
	remainder := (qtdPosts % limit)
	if remainder == 0 {
		return total, nil
	}

	return total + 1, nil
}
