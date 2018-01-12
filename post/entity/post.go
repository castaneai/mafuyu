package entity

import "time"

type Post struct {
	ID           int64      `json:"id" datastore:"-" boom:"id"`
	Title        string     `json:"title"`
	Tags         []string   `json:"tags"`
	Pages        []PostPage `json:"pages"`
	ThumbnailURL string     `json:"thumbnail_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
