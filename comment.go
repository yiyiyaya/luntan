package main

import (
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

//curl  -d '{"Content":"lsm"}' -X POST localhost:8002/posts/zls7e6j4dyhyo/comment
func CreateComment(c *gin.Context) {
	var comment Comment

	err := c.BindJSON(&comment)
	if err != nil {
		log.Print("err111111111111111", err)
		c.JSON(http.StatusBadRequest, Error{
			Message: err.Error(),
		})
		return
	}

	if len(comment.Content) > 10000 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "content too long",
		})
		return
	}
	if comment.Content == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "Content can't nil",
		})
		return
	}
	postID := c.Param("postid")
	if postID == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "postID can't nil",
		})
		return
	}
	if len(postID) > 32 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "postID too long",
		})
		return
	}
	commentID, err := db.CreateComment(postID, comment.Content)
	if err != nil {
		log.Print("err22222222222222222", err)
		c.JSON(http.StatusInternalServerError, Error{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		Id int64
	}{
		Id: commentID,
	})
}

//delete
//curl   -X DELETE localhost:8002/posts/737x7p67553tw/comments/1
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentid")
	if commentID == "" {
		log.Print("commentid can't 0")
		return
	}
	if len(commentID) > 20 {
		log.Print("commentid too long")
		return
	}
	id, err := strconv.Atoi(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{
			Message: err.Error(),
		})
		return
	}
	err = db.DeleteComment(id)
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

//curl  -d '{"Content":"lsm"}' -X PATCH localhost:8002/posts/zls7e6j4dyhyo/comments/4
func UpdateComment(c *gin.Context) {
	var comment Comment
	err := c.BindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, Error{
			Message: err.Error(),
		})
		return
	}
	if len(comment.Content) > 10000 {
		c.JSON(http.StatusBadRequest, Error{
			Message: "content too long",
		})
		return
	}
	if comment.Content == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "Content can't nil",
		})
		return
	}
	id := c.Param("commentid")
	if id == "" {
		c.JSON(http.StatusBadRequest, Error{
			Message: "id can't nil",
		})
		return
	}
	err = db.UpdateComment(id, comment.Content)
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
