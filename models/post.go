package models

type Post struct {
    Username  string
    Email string
    Content string
}

func NewPost(username, email, content string) *Post {
    return &Post{username, email, content}
}