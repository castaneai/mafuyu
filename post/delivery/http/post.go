package http

import (
	"context"
	"fmt"
	"github.com/castaneai/mafuyu/post/entity"
	"github.com/castaneai/mafuyu/post/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
)

func getPostRepo(ctx context.Context) (repository.PostRepository, error) {
	return repository.NewDatastorePostRepository(ctx)
}

func validatePost(post *entity.Post) error {
	if len(post.Pages) < 1 {
		return errors.New("post.Pages len must be > 0")
	}
	return nil
}

func handleError(ctx context.Context, gc *gin.Context, err error, message string) {
	log.Errorf(ctx, "%+v", err)
	gc.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "ERROR", "message": message})
}

func insertPost(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		handleError(ctx, c, err, "failed to create post repository")
		return
	}

	post := &entity.Post{}
	if err := c.BindJSON(post); err != nil {
		handleError(ctx, c, err, "failed to bind json to post entity")
		return
	}
	if err := validatePost(post); err != nil {
		handleError(ctx, c, err, fmt.Sprintf("%s", err))
		return
	}

	duplicationFound := false
	status := "OK"
	for _, source := range post.Sources {
		dpost, err := repo.FindBySourceID(source.ID)
		if err != nil {
			handleError(ctx, c, err, fmt.Sprintf("failed to find post by source id: %+v", err))
			return
		}
		if dpost != nil {
			duplicationFound = true
			status = "DUPLICATION FOUND"
			post = dpost
			break
		}
	}

	if !duplicationFound {
		post, err = repo.Insert(post)
		if err != nil {
			handleError(ctx, c, err, fmt.Sprintf("failed to insert post"))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"posts":  []entity.Post{*post},
	})
}

func searchPost(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		handleError(ctx, c, err, "failed to create post repository")
		return
	}

	posts, err := repo.Search(c.Query("q"))
	if err != nil {
		handleError(ctx, c, err, "failed to search posts")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"posts":  posts,
	})
}

func getPost(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		handleError(ctx, c, err, "failed to create post repository")
		return
	}

	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(ctx, c, err, "failed to convert id param to int")
		return
	}

	post, err := repo.Find(int64(pid))
	if err != nil {
		handleError(ctx, c, err, fmt.Sprintf("failed to find post with id: %d", pid))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"post":   post,
	})
}

func searchTag(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		handleError(ctx, c, err, "failed to create post repository")
		return
	}

	tags, err := repo.SearchTag(c.Query("q"))
	if err != nil {
		handleError(ctx, c, err, "failed to search tags")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"tags":   tags,
	})
}

func countPost(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	repo, err := getPostRepo(ctx)
	if err != nil {
		handleError(ctx, c, err, "failed to create post repository")
		return
	}

	count, err := repo.Count(c.Query("q"))
	if err != nil {
		handleError(ctx, c, err, "failed to count posts")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"count":  count,
	})
}

func Init() *gin.Engine {
	r := gin.Default()

	// cors
	r.Use(cors.Default())

	r.GET("/post", searchPost)
	r.GET("/post/:id", getPost)
	r.POST("/post", insertPost)
	r.GET("/tag", searchTag)
	r.GET("/count", countPost)
	return r
}
