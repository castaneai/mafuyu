package usecase

import (
	"context"
	"errors"
	"github.com/castaneai/mafuyu/post/entity"
	"github.com/castaneai/mafuyu/post/repository"
)

type postUsecase struct {
	repo repository.PostRepository
}

func NewPostUsecase(repo repository.PostRepository) PostUsecase {
	return &postUsecase{
		repo: repo,
	}
}

func (uc *postUsecase) Search(ctx context.Context, keyword string) ([]*entity.Post, error) {
	return uc.repo.Search(keyword)
}

func (uc *postUsecase) Find(ctx context.Context, postID int64) (*entity.Post, error) {
	return uc.repo.Find(postID)
}

func (uc *postUsecase) Insert(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if err := validatePost(uc.repo, post); err != nil {
		return nil, err
	}
	return uc.repo.Insert(post)
}

func (uc *postUsecase) SearchTags(ctx context.Context, keyword string) ([]*repository.TagInfo, error) {
	return uc.repo.SearchTag(keyword)
}

func (uc *postUsecase) Count(ctx context.Context, keyword string) (int, error) {
	return uc.repo.Count(keyword)
}

func validatePost(repo repository.PostRepository, post *entity.Post) error {
	if len(post.Pages) < 1 {
		return errors.New("post.Pages len must be > 0")
	}

	for _, source := range post.Sources {
		dpost, err := repo.FindBySourceID(source.ID)
		if err != nil {
			return err
		}
		if dpost != nil {
			return errors.New("duplication post found")
		}
	}

	return nil
}
