package repository

import "github.com/castaneai/mafuyu/post/entity"

type TagInfo struct {
	Tag       string `json:"tag"`
	PostCount int    `json:"post_count"`
}

type PostRepository interface {
	Find(id int64) (*entity.Post, error)
	Search(keyword string) ([]*entity.Post, error)
	SearchTag(keyword string) ([]*TagInfo, error)
	Insert(post *entity.Post) (*entity.Post, error)
	Update(post *entity.Post) (*entity.Post, error)
	Delete(post *entity.Post) error
}
