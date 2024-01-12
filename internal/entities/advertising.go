package entities

import (
	"time"
)

type AdvertisingRequest struct {
	Status   string      `json:"status" validate:"required,max=10"`
	Category string      `json:"category" validate:"required,max=10"`
	Post     PostRequest `json:"post" validate:"required"`
}

type Advertising struct {
	Id       int
	Status   string
	Category string
	UserId   int
	Post     Post
	PostId   int
}

func (a *Advertising) NewPost(post PostRequest) error {
	parseddate, err := time.Parse("02-01-2006", post.Date)
	if err != nil {
		return err
	}
	a.Post = Post{
		Title:       post.Title,
		Description: post.Description,
		Date:        parseddate,
	}
	return nil
}
