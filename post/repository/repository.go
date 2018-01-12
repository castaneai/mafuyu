package repository

import "github.com/castaneai/mafuyu/post/entity"

type PostRepository interface {
	Find(id int64) (*entity.Post, error)
	Search(keyword string) ([]*entity.Post, error)
	Insert(post *entity.Post) (*entity.Post, error)
	Update(post *entity.Post) (*entity.Post, error)
}
