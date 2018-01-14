package main

import (
	"github.com/castaneai/mafuyu/post/delivery/http"
	_ "go.mercari.io/datastore/clouddatastore"
)

func main() {
	gin := http.Init()
	gin.Run("127.0.0.1:8080")
}
