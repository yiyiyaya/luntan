package main

import "time"

type Post struct {
	Id         string
	Title      string
	CreateTime time.Time
}
type Comment struct {
	Id         int64
	Content    string
	CreateTime time.Time
}
