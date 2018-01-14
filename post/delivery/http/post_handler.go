// TODO: appengine 依存を取り払う
package http

import (
	"context"
	"fmt"
	"github.com/castaneai/mafuyu/post/entity"
	"github.com/castaneai/mafuyu/post/repository"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mercari.io/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

const (
	ProjectID = "morning-tide"
)

func getPostRepo(ctx context.Context) (repository.PostRepository, error) {
	opts := datastore.WithProjectID(ProjectID)
	return repository.NewDatastorePostRepository(ctx, opts)
}

func getMockPost() (*entity.Post, error) {
	return &entity.Post{
		Title:        "12345",
		Tags:         []string{"tag1", "tag2"},
		Pages:        []entity.PostPage{{ContentURL: "https://cs.sankakucomplex.com/data/sample/da/ae/sample-daaeb9dc4bf27b74264eb4c8a1d2b15e.jpg?e=1516033173&m=voJk6GLXuwBDDxBBo8teVQ"}},
		ThumbnailURL: "https://cs.sankakucomplex.com/data/preview/af/56/af565d163dd29df40e67884dbc000b0e.jpg",
	}, nil
}

func insertPost(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "failed to create post repo"})
		return
	}

	post, err := getMockPost()
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": fmt.Sprintf("failed to get sankaku post: %v", err)})
		return
	}

	post, err = repo.Insert(post)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": fmt.Sprintf("failed to insert post: %+v", errors.WithStack(err))})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "posts": []entity.Post{*post}})
}

func searchPost(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "failed to create post repo"})
		return
	}
	posts, err := repo.Search(c.Param("q"))
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "failed to search"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func Init() *gin.Engine {
	r := gin.Default()
	postAPI := r.Group("/api/v1/post")
	postAPI.GET("/", searchPost)
	postAPI.POST("/", insertPost)
	return r
}
