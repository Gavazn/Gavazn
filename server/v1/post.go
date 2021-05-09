package v1

import (
	"net/http"

	"github.com/Gavazn/Gavazn/internal/post"
	"github.com/Gavazn/Gavazn/internal/user"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type postForm struct {
	Title      string               `json:"title" form:"title"`
	Content    string               `json:"content" form:"content"`
	Categories []primitive.ObjectID `json:"categories" form:"categories"`
	Tags       []string             `json:"tags" form:"tags"`
	Thumbnail  primitive.ObjectID   `json:"thumbnail" form:"thumbnail"`
}

/**
 * @api {post} /api/v1/posts add post
 * @apiVersion 1.0.0
 * @apiName addPost
 * @apiGroup Post
 *
 * @apiParam {String} title post title
 * @apiParam {String} content post content
 * @apiParam {String[]} categories list of category id
 * @apiParam {String[]} tags list of tag
 * @apiParam {String} thumbnail thumbnail id
 *
 * @apiSuccess {String} message success message.
 * @apiSuccess {Object} post post model
 *
 * @apiError {String} error api error message
 */
func addPost(ctx echo.Context) error {
	u := ctx.Get("user").(*user.User)

	form := new(postForm)
	if err := ctx.Bind(form); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	p := &post.Post{
		User:       u.ID,
		Title:      form.Title,
		Content:    form.Content,
		Categories: form.Categories,
		Tags:       form.Tags,
		Thumbnail:  form.Thumbnail,
	}

	if err := p.Save(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "post created successfully",
		"post":    p,
	})
}

/**
 * @api {put} /api/v1/posts edit post
 * @apiVersion 1.0.0
 * @apiName editPost
 * @apiGroup Post
 *
 * @apiParam {String} title post title
 * @apiParam {String} content post content
 * @apiParam {String[]} categories list of category id
 * @apiParam {String[]} tags list of tag
 * @apiParam {String} thumbnail thumbnail id
 *
 * @apiSuccess {String} message success message.
 * @apiSuccess {Object} post post model
 *
 * @apiError {String} error api error message
 */
func editPost(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	form := new(postForm)
	if err := ctx.Bind(form); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	p, err := post.FindOne(bson.M{"_id": id})
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	p.Title = form.Title
	p.Content = form.Content
	p.Categories = form.Categories
	p.Tags = form.Tags
	p.Thumbnail = form.Thumbnail

	if err := p.Save(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "post updated successfully",
		"post":    p,
	})
}
