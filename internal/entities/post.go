package entities

import "time"

type Post struct {
	Id          int
	Title       string
	Description string
	Date        *time.Time
	Pet_id      int
}
