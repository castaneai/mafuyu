package http

import (
	"context"
	"github.com/castaneai/mafuyu/post/entity"
	"github.com/castaneai/mafuyu/post/repository"
	"github.com/castaneai/mafuyu/post/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
)

func handleError(ctx context.Context, gc *gin.Context, err error, message string) {
	log.Errorf(ctx, "%+v", err)
	gc.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "ERROR", "message": message})
}

func insertPost(ctx context.Context, ginc *gin.Context, uc usecase.PostUsecase) (gin.H, error) {
	post := &entity.Post{}
	if err := ginc.BindJSON(post); err != nil {
		return nil, err
	}

	post, err := uc.Insert(ctx, post)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"posts": []*entity.Post{post},
	}, nil
}

func searchPost(ctx context.Context, ginc *gin.Context, uc usecase.PostUsecase) (gin.H, error) {
	posts, err := uc.Search(ctx, ginc.Query("q"))
	if err != nil {
		return nil, err
	}

	return gin.H{
		"posts": posts,
	}, nil
}

func getPost(ctx context.Context, ginc *gin.Context, uc usecase.PostUsecase) (gin.H, error) {
	postId, err := strconv.Atoi(ginc.Param("id"))
	if err != nil {
		return nil, err
	}

	post, err := uc.Find(ctx, int64(postId))
	if err != nil {
		return nil, err
	}

	return gin.H{
		"post": post,
	}, nil
}

func searchTag(ctx context.Context, ginc *gin.Context, uc usecase.PostUsecase) (gin.H, error) {
	tags, err := uc.SearchTags(ctx, ginc.Query("q"))
	if err != nil {
		return nil, err
	}

	return gin.H{
		"tags": tags,
	}, nil
}

func countPost(ctx context.Context, ginc *gin.Context, uc usecase.PostUsecase) (gin.H, error) {
	count, err := uc.Count(ctx, ginc.Query("q"))
	if err != nil {
		return nil, err
	}

	return gin.H{
		"count": count,
	}, nil
}

func wrapHandler(innerHandler func(ctx context.Context, gin *gin.Context, uc usecase.PostUsecase) (gin.H, error)) gin.HandlerFunc {
	// wrapping common request handling and create usecase
	return func(ginc *gin.Context) {
		ctx := appengine.NewContext(ginc.Request)
		repo, err := repository.NewDatastorePostRepository(ctx)
		if err != nil {
			handleError(ctx, ginc, err, "failed to create post repository")
			return
		}
		uc := usecase.NewPostUsecase(repo)

		res, err := innerHandler(ctx, ginc, uc)

		if err != nil {
			handleError(ctx, ginc, err, "internal server error")
			return
		}
		ginc.JSON(http.StatusOK, res)
	}
}

func Init() *gin.Engine {
	r := gin.Default()

	// cors
	r.Use(cors.Default())

	r.GET("/post", wrapHandler(searchPost))
	r.GET("/post/:id", wrapHandler(getPost))
	r.POST("/post", wrapHandler(insertPost))
	r.GET("/tag", wrapHandler(searchTag))
	r.GET("/count", wrapHandler(countPost))
	return r
}
