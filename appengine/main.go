package main

import (
	mafuyuHttp "github.com/castaneai/mafuyu/post/delivery/http"
	_ "go.mercari.io/datastore/aedatastore"
	"net/http"
)

func init() {
	h := mafuyuHttp.Init()
	http.Handle("/", h)
}
