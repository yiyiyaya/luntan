package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string
}

var db DB

func main() {

	//
	var err error
	db, err = NewMysql()
	if err != nil {
		log.Fatal(err)
	}

	//
	router := gin.Default()
	router.GET("/", index)
	router.POST("/post", CreatePost)
	router.POST("/post/comment", CreateComment)

	router.Run(":9091")

}
func index(c *gin.Context) {
	c.String(http.StatusOK, "hello world")

}
