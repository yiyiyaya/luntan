package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PostCreate struct {
	Post    Post
	Comment Comment
}

//curl  -d '{"Post":{"Title":"lsm"},"Comment":{"Content":"fdhgdfhfghfgh"}}' -X POST localhost:8002/posts
func CreatePost(c *gin.Context) {
	var post PostCreate
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{
			Message: err.Error(),
		})
		return

	}
	if post.Post.Title == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "title is blank",
		})
		return
	}
	if len(post.Post.Title) > 200 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "title too long",
		})
		return
	}
	if post.Comment.Content == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "Comment is blank",
		})
		return
	}
	if len(post.Comment.Content) > 10000 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "title too long",
		})
		return
	}
	postID := generateID()
	err = db.CreatePost(postID, post.Post.Title, post.Comment.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		Id string
	}{
		Id: postID,
	})
}

func getPosts(c *gin.Context) ([]Post, int, error) {
	var err error
	var fromTime time.Time
	inputTime, ok := c.GetQuery("fromtime")
	if ok {
		var err error
		fromTime, err = time.Parse("20060102150405", inputTime)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	} else {
		fromTime = time.Now()
	}
	inputCount, ok := c.GetQuery("count")
	var count int
	if ok {
		var err error
		count, err = strconv.Atoi(inputCount)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	} else {
		count = 5
	}
	posts, err := db.ListPostsByCreateTime(fromTime, count)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return posts, http.StatusOK, nil
}

//curl   -X GET localhost:8002/posts?fromtime=20170101160000&count=12
func ListPost(c *gin.Context) {
	posts, code, err := getPosts(c)
	if err != nil {
		c.JSON(code, Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, struct {
		Posts []Post
	}{
		Posts: posts,
	})

}

//curl   -X GET localhost:8002/posts/rdcoe4jydshmo/comment
func GetPost(c *gin.Context) {
	postID := c.Param("postid")
	if postID == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "postID can't be blank",
		})
		return
	}
	if len(postID) > 32 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "postID too long",
		})
		return
	}
	post, comments, err := db.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		Post    Post
		Comment []Comment
	}{
		Post:    post,
		Comment: comments,
	})

}

//curl   -X DELETE localhost:8002/posts/737x7p67553tw
func DeletePost(c *gin.Context) {
	postID := c.Param("postid")
	if postID == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "postID can't be blank",
		})
		return
	}
	if len(postID) > 32 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "postID too long",
		})
		return
	}
	err := db.DeletePost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		Message string
	}{
		Message: "success",
	})
}

//curl  -d '{"Post":{"Title":"123"},"Comment":{"Content":"456"}}' -X PATCH localhost:8002/posts/zls7e6j4dyhyo
func UpdatePost(c *gin.Context) {
	var post PostCreate
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{
			Message: err.Error(),
		})
		return

	}
	if post.Post.Title == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "title is blank",
		})
		return
	}
	if len(post.Post.Title) > 200 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "title too long",
		})
		return
	}
	if post.Comment.Content == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "Comment is blank",
		})
		return
	}
	if len(post.Comment.Content) > 10000 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "Content too long",
		})
		return
	}
	postID := c.Param("postid")
	err = db.UpdatePost(postID, post.Post.Title, post.Comment.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		Message string
	}{
		Message: "update success",
	})
}
