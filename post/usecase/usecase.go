package usecase

import (
	"context"
	"github.com/castaneai/mafuyu/post/entity"
	"github.com/castaneai/mafuyu/post/repository"
)

type PostUsecase interface {
	Search(ctx context.Context, keyword string) ([]*entity.Post, error)
	Find(ctx context.Context, postID int64) (*entity.Post, error)
	Insert(ctx context.Context, post *entity.Post) (*entity.Post, error)
	SearchTags(ctx context.Context, keyword string) ([]*repository.TagInfo, error)
	Count(ctx context.Context, keyword string) (int, error)
}
