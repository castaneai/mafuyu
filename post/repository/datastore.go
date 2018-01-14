package repository

import (
	"context"
	"github.com/castaneai/mafuyu/post/entity"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"strings"
	"time"
)

const (
	searchLimit = 100
)

type datastorePostRepository struct {
	boom *boom.Boom
}

func (repo *datastorePostRepository) Find(id int64) (*entity.Post, error) {
	data := &entity.Post{ID: id}
	if err := repo.boom.Get(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *datastorePostRepository) Search(keyword string) ([]*entity.Post, error) {
	// TODO: keyword parser
	q := repo.boom.NewQuery("Post")
	for tag := range strings.Split(keyword, " ") {
		q = q.Filter("tags =", tag)
	}
	// 念のため上限かけておく
	q = q.Limit(searchLimit)
	var posts []*entity.Post
	if _, err := repo.boom.GetAll(q, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *datastorePostRepository) Insert(post *entity.Post) (*entity.Post, error) {
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	if _, err := repo.boom.Put(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *datastorePostRepository) Update(post *entity.Post) (*entity.Post, error) {
	// datastore では insert/update 共に put であり、キーが一致すれば自動的に上書きとなる
	// https://cloud.google.com/appengine/docs/standard/go/datastore/creating-entities#Go_Updating_entities
	return repo.Insert(post)
}

func (repo *datastorePostRepository) Delete(post *entity.Post) error {
	if err := repo.boom.Delete(post); err != nil {
		return err
	}
	return nil
}

func NewDatastorePostRepository(ctx context.Context, opts datastore.ClientOption) (PostRepository, error) {
	ds, err := datastore.FromContext(ctx, opts)
	if err != nil {
		return nil, err
	}
	b := boom.FromClient(ctx, ds)
	repo := &datastorePostRepository{boom: b}
	return repo, nil
}
