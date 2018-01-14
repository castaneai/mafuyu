package converter

import (
	"github.com/castaneai/mafuyu/post/entity"
	"github.com/castaneai/sankaku"
)

type SankakuConverter struct{}

func NewSankakuConverter() (*SankakuConverter, error) {
	return &SankakuConverter{}, nil
}

func (s *SankakuConverter) Convert(sp *sankaku.Post, sd *sankaku.PostDetail) (*entity.Post, error) {
	return &entity.Post{
		Title:        sp.ID,
		Tags:         sp.Tags,
		Pages:        []entity.PostPage{{ContentURL: sd.OriginalURL}},
		ThumbnailURL: sp.ThumbnailURL,
	}, nil
}
