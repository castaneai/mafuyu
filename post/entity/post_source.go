package entity

type PostSource struct {
	ID    string `json:"id" datastore:"id"`
	Title string `json:"title" datastore:"title,noindex"`
	URL   string `json:"url" datastore:"url,noindex"`
}
