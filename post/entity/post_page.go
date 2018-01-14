package entity

type PostPage struct {
	ContentURL string `json:"content_url" datastore:"content_url,noindex"`
}
