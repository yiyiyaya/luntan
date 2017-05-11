package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const CreatePostSql = `CREATE  TABLE IF NOT EXISTS POST
(
	ID         VARCHAR(32) primary key,
	TITLE      VARCHAR(256),
	CREATETIME DATETIME
)DEFAULT CHARSET=UTF8;
`
const CreateCommentSql = `CREATE TABLE IF NOT EXISTS COMMENT
(
	ID         BIGINT NOT NULL AUTO_INCREMENT,
	CONTENT    VARCHAR(10000) NOT NULL,
	CREATETIME DATETIME,
	PRIMARY KEY (ID)
)DEFAULT CHARSET=UTF8;`

type MySql struct {
	db *sql.DB
}

func (ms *MySql) CreatePost(postID, title, content string) (err error) {
	tx, err := ms.db.Begin()
	if err != nil {
		return
	}
	now := time.Now().UTC().Format("2006-01-02 15:04:05.999999")
	_, err = tx.Exec("insert into POST  (ID, TITLE, CREATETIME) values (?, ?, ?)", postID, title, now)
	if err != nil {
		tx.Rollback()
		log.Println(errors.New("post insert err"))
		return
	}
	_, err = tx.Exec("insert into COMMENT  (CONTENT, CREATETIME) values (?, ?)", content, now)
	if err != nil {
		tx.Rollback()
		log.Println(errors.New("comment insert err"))
		return
	}
	//id, err := result.LastInsertId()
	err = tx.Commit()
	return

}
func (ms *MySql) CreateComment(postID, content string) (commentID int64, err error) {

	return 0, errors.New("not implemented yet")
}

func NewMysql() (mysql *MySql, err error) {
	db_addr := os.Getenv("MYSQL_ADDR")
	db_port := os.Getenv("MYSQL_PORT")
	db_database := os.Getenv("MYSQL_DATABASE")
	db_user := os.Getenv("MYSQL_USER")
	db_password := os.Getenv("MYSQL_PASSWORD")
	db_url := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true`,
		db_user, db_password, db_addr, db_port, db_database)
	log.Printf("connect to %v", db_url)
	db, err := sql.Open("mysql", db_url)
	if err != nil {
		log.Println("err:", err)
		return
	}
	//ping一个服务器链接，如果没有链接则要重新链接
	err = db.Ping()
	if err != nil {
		return
	}
	_, err = db.Exec(CreatePostSql)
	if err != nil {
		log.Println("CreatePostSQl:", err)
		return
	}
	_, err = db.Exec(CreateCommentSql)
	if err != nil {
		log.Println("CreateCommentSQl:", err)
		return
	}
	mysql = &MySql{
		db: db,
	}
	return mysql, nil
}
