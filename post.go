package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostCreate struct {
	Post    Post
	Comment Comment
}

//curl  -d '{"Post":{"Title":"lsm"},"Comment":{"Content":"fdhgdfhfghfgh"}}' -X POST localhost:8002/post
func CreatePost(c *gin.Context) {
	var post PostCreate
	if c.BindJSON(&post) == nil {
		c.JSON(200, gin.H{"received": &post})
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
	err := db.CreatePost(postID, post.Post.Title, post.Comment.Content)
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
