package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Postgre struct {
	db *sql.DB
}

func (pg *Postgre) CreatePost(title, content string) (postID string, err error) {

	return "", nil
}
func (pg *Postgre) CreateComment(postID, content string) (commentID int64, err error) {

	return 0, nil
}

func NewPostgre() (pg *Postgre, err error) {
	DB_ADDR := os.Getenv("MYSQL_ADDR")
	DB_PORT := os.Getenv("MYSQL_PORT")
	DB_DATABASE := os.Getenv("MYSQL_DATABASE")
	DB_USER := os.Getenv("MYSQL_USER")
	DB_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	DB_URL := fmt.Sprintf(`%S:%s@tcp(%S:%S)/%S?charset=utf-8&parseTime=true`,
		DB_USER, DB_PASSWORD, DB_ADDR, DB_PORT, DB_DATABASE)
	log.Printf("connect to %v", DB_URL)
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Println("err:", err)
		return
	}
	//ping一个服务器链接，如果没有链接则要重新链接
	err = db.Ping()
	if err != nil {
		log.Println("err:", err)
		return
	}
	pg = &Postgre{
		db: db,
	}
	return pg, nil
}
