package services

import (
	"context"

	"github.com/devGabrielb/AmiFind/internal/repositories"
)

type PostService interface {
	Get(ctx context.Context, page int) (*repositories.PostPaged, error)
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo: repo}
}

func (ps *postService) Get(ctx context.Context, page int) (*repositories.PostPaged, error) {
	postsPaged, err := ps.repo.GetPostPaged(ctx, page)
	if err != nil {
		return nil, err
	}

	return postsPaged, nil

}
