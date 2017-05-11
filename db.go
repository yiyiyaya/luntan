package main

import (
	"encoding/base32"
	"time"
)

type DB interface {
	CreatePost(postID, title, content string) (err error)
	CreateComment(postID, content string) (commentID int64, err error)
}

var base32Encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

func generateID() string {
	t := time.Now().UnixNano()
	bs := make([]byte, 8)
	for i := uint(0); i < 8; i++ {
		bs[i] = byte((t >> i) & 0xff)
	}

	dest := make([]byte, 16)
	base32Encoding.Encode(dest, bs)
	return string(dest[:13])
}
