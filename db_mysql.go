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
	POSTID     VARCHAR(32), 
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
	_, err = tx.Exec("insert into COMMENT  (CONTENT, CREATETIME, POSTID) values (?, ?, ?)", content, now, postID)
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
	createtime := time.Now().UTC().Format("2006-01-02 15:04:05.999999")
	rs, err := ms.db.Exec("insert into COMMENT  (CONTENT, CREATETIME, POSTID) values (?, ?, ?)", content, createtime, postID)
	if err != nil {
		log.Print("insert commentcreate err")
		return
	}
	commentID, err = rs.LastInsertId()
	return
}
func (ms *MySql) ListPostsByCreateTime(startTime time.Time, maxReturns int) (posts []Post, err error) {
	rows, err := ms.db.Query("select CREATETIME, ID, TITLE from POST where CREATETIME < ? limit ?", startTime.UTC().Format("2006-01-02 15:04:05.999999"), maxReturns)
	if err != nil {
		log.Print("select faild")
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := Post{}
		err = rows.Scan(
			&p.CreateTime,
			&p.Id,
			&p.Title,
		)
		if err != nil {
			log.Print("scan err")
			return
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		log.Print("rows err")
		return
	}
	return
}
func (ms *MySql) GetPost(postID string) (post Post, comments []Comment, err error) {
	row := ms.db.QueryRow("select CREATETIME, ID,  TITLE from POST where ID = ?", postID)
	err = row.Scan(
		&post.CreateTime,
		&post.Id,
		&post.Title,
	)
	if err != nil {
		log.Print("scan post err")
		return
	}

	rows, err := ms.db.Query("select CREATETIME, ID, CONTENT from COMMENT where POSTID = ?", postID)
	if err != nil {
		log.Print("select comment err")
		return
	}
	for rows.Next() {
		c := Comment{}
		err = rows.Scan(
			&c.CreateTime,
			&c.Id,
			&c.Content,
		)
		if err != nil {
			log.Print("scan comment err")
			return
		}
		if err = rows.Err(); err != nil {
			return
		}
		comments = append(comments, c)
	}
	return
}

//delete
func (ms *MySql) DeletePost(postID string) (err error) {
	_, err = ms.db.Exec("delete from POST where ID = ?", postID)
	if err != nil {
		log.Print("delete err")
		return
	}
	log.Print("delete success")
	return
}
func (ms *MySql) DeleteComment(commentID int) (err error) {
	_, err = ms.db.Exec("delete from COMMENT where ID = ?", commentID)
	if err != nil {
		log.Print("delete err")
		return
	}
	return
}

//update

func (ms *MySql) UpdatePost(postID, title, content string) (err error) {
	tx, err := ms.db.Begin()
	if err != nil {
		return
	}
	now := time.Now().UTC().Format("2006-01-02 15:04:05.999999")
	_, err = ms.db.Exec("update POST set  TITLE = ?, CREATETIME = ? where ID = ?", title, now, postID)
	if err != nil {
		log.Print("update post err")
		tx.Rollback()
		return
	}
	_, err = ms.db.Exec("update COMMENT set CONTENT = ?, CREATETIME = ? where POSTID = ?", content, now, postID)
	if err != nil {
		log.Print("update comment err")
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}
func (ms *MySql) UpdateComment(id, content string) (err error) {
	createtime := time.Now().UTC().Format("2006-01-02 15:04:05.999999")
	_, err = ms.db.Exec("update COMMENT set CONTENT = ?, CREATETIME = ? where ID = ?", content, createtime, id)
	if err != nil {
		log.Print("update Comment  err")
		return
	}
	return

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
