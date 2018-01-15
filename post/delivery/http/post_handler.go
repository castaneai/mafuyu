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

func validatePost(post *entity.Post) error {
	if len(post.Pages) < 1 {
		return errors.New("post.Pages len must be > 0")
	}
	return nil
}

func insertPost(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "failed to create post repo"})
		return
	}

	post := &entity.Post{}
	if err := c.BindJSON(post); err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "failed to bind request JSON"})
		return
	}
	if err := validatePost(post); err != nil {
		log.Errorf(ctx, "%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": fmt.Sprintf("%s", err)})
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
		"status": "OK",
		"posts":  posts,
	})
}

func Init() *gin.Engine {
	r := gin.Default()
	postAPI := r.Group("/api/v1/post")
	postAPI.GET("/", searchPost)
	postAPI.POST("/", insertPost)
	return r
}
