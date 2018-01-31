package repository

import (
	"context"
	"github.com/castaneai/mafuyu/post/entity"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	kindName    = "Post"
	searchLimit = 10000
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
	q := repo.boom.NewQuery(kindName)
	for _, tag := range strings.Split(keyword, " ") {
		if tag != "" {
			q = q.Filter("tags =", tag)
		}
	}
	// 念のため上限かけておく
	q = q.Limit(searchLimit)
	q = q.Order("-created_at")
	var posts []*entity.Post
	if _, err := repo.boom.GetAll(q, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *datastorePostRepository) Insert(post *entity.Post) (*entity.Post, error) {
	now := time.Now()
	post.CreatedAt = now
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

func (repo *datastorePostRepository) SearchTag(keyword string) ([]*TagInfo, error) {
	q := repo.boom.NewQuery(kindName).Filter("tags >=", keyword).Filter("tags <", keyword+string([]rune{utf8.MaxRune}))
	q = q.Limit(searchLimit)
	var posts []*entity.Post
	if _, err := repo.boom.GetAll(q, &posts); err != nil {
		return nil, err
	}

	tagCountMap := make(map[string]int)
	for _, post := range posts {
		for _, tag := range post.Tags {
			if !strings.HasPrefix(tag, keyword) {
				continue
			}
			if _, ok := tagCountMap[tag]; !ok {
				tagCountMap[tag] = 0
			}
			tagCountMap[tag] += 1
		}
	}

	tagInfos := make([]*TagInfo, len(tagCountMap))
	i := 0
	for tag, count := range tagCountMap {
		tagInfos[i] = &TagInfo{
			Tag:       tag,
			PostCount: count,
		}
		i++
	}

	sort.Slice(tagInfos, func(i, j int) bool {
		return tagInfos[i].PostCount > tagInfos[j].PostCount
	})

	return tagInfos, nil
}

func (repo *datastorePostRepository) Count(keyword string) (int, error) {
	// TODO: keyword parser
	q := repo.boom.NewQuery(kindName)
	for _, tag := range strings.Split(keyword, " ") {
		if tag != "" {
			q = q.Filter("tags =", tag)
		}
	}
	count, err := repo.boom.Count(q)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *datastorePostRepository) FindBySourceID(sourceID string) (*entity.Post, error) {
	q := repo.boom.NewQuery(kindName).Filter("sources.id =", sourceID).Limit(1)
	var posts []*entity.Post
	if _, err := repo.boom.GetAll(q, &posts); err != nil {
		return nil, err
	}
	if len(posts) < 1 {
		return nil, nil
	}
	return posts[0], nil
}

func NewDatastorePostRepository(ctx context.Context, opts ...datastore.ClientOption) (PostRepository, error) {
	ds, err := datastore.FromContext(ctx, opts...)
	if err != nil {
		return nil, err
	}
	b := boom.FromClient(ctx, ds)
	repo := &datastorePostRepository{boom: b}
	return repo, nil
}
