package main

import (
	"log"

	"html/template"

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
	// router.GET("/", index)
	//router.Use()
	router.GET("/", Home)
	//router.POST("/posts", pageCreatePost)
	router.GET("/posts/:postid", ShowPost)
	//	router.Post("/posts/:postid", pageCreateComment)
	router.Static("/static", "./static")
	router.POST("/api/v1/posts", CreatePost)
	router.POST("/api/v1/posts/:postid/comment", CreateComment)
	router.GET("/api/v1/posts", ListPost)
	router.GET("/api/v1/posts/:postid/comment", GetPost)
	router.DELETE("/api/v1/posts/:postid", DeletePost)
	router.DELETE("/api/v1/posts/:postid/comments/:commentid", DeleteComment)
	router.PATCH("/api/v1/posts/:postid", UpdatePost)
	router.PATCH("/api/v1/posts/:postid/comments/:commentid", UpdateComment)

	router.Run(":9091")

}

//func index(c *gin.Context) {
//	c.String(http.StatusOK, "hello world")
//}
var isProduction = false

var homeTemplate = template.Must(template.ParseFiles("./template/index.html"))

func LoadHomeTemplate() *template.Template {
	if isProduction {
		return homeTemplate
	}
	return template.Must(template.ParseFiles("./template/index.html"))

}

func Home(c *gin.Context) {
	posts, _, err := getPosts(c)
	if err != nil {
		LoadHomeTemplate().Execute(c.Writer, struct {
			Posts []Post
			Title string
			Error string
		}{
			Posts: posts,
			Title: "Post List",
			Error: "Error:" + err.Error(),
		})
		return
	}

	LoadHomeTemplate().Execute(c.Writer, struct {
		Posts []Post
		Title string
		Error string
	}{
		Posts: posts,
		Title: "Post List",
		Error: "",
	})
}

var postTemplate = template.Must(
	template.New("hello").Funcs(
		template.FuncMap{
			"IsEven": IsEven,
		},
	).ParseFiles("./template/post.html"),
)

func LoadPostTemplate() *template.Template {
	if isProduction {
		return postTemplate
	}
	return template.Must(
		template.New("hello").Funcs(
			template.FuncMap{
				"IsEven": IsEven,
			},
		).ParseFiles("./template/post.html"),
	)

}
func ShowPost(c *gin.Context) {
	postID := c.Param("postid")
	post, comments, err := db.GetPost(postID)
	if err != nil {
		LoadPostTemplate().ExecuteTemplate(c.Writer, "post.html", struct {
			Post     Post
			Comments []Comment
			Title    string
			Error    string
		}{
			Post:     post,
			Comments: nil,
			Title:    "Post",
			Error:    err.Error(),
		})
		return
	}
	err = LoadPostTemplate().ExecuteTemplate(c.Writer, "post.html", struct {
		Post     Post
		Comments []Comment
		Title    string
		Error    string
	}{
		Post:     post,
		Comments: comments,
		Title:    "Post",
		Error:    "",
	})
	log.Println(err)
	return
}
