package entity

import (
	"encoding/json"
	"time"
)

type Post struct {
	ID           int64      `json:"id" datastore:"-" boom:"id"`
	Title        string     `json:"title" datastore:"title,noindex"`
	Tags         []string   `json:"tags" datastore:"tags"`
	Pages        []PostPage `json:"pages" datastore:"pages,noindex,flatten"` // flatten 忘れずに！！ https://qiita.com/vvakame/items/9310bcb5a4e87888d505#%E7%A7%BB%E8%A1%8C%E3%81%AE%E6%B3%A8%E6%84%8F%E7%82%B9
	ThumbnailURL string     `json:"thumbnail_url" datastore:"thumbnail_url,noindex"`
	CreatedAt    time.Time  `json:"created_at" datastore:"created_at,noindex"`
	UpdatedAt    time.Time  `json:"updated_at" datastore:"updated_at"`
}

func (p *Post) String() string {
	j, _ := json.Marshal(p)
	return string(j)
}
