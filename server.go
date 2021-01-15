package main

import (
	"fmt"
	"net/http"
	"github.com/thedevsaddam/renderer"
	//"encoding/json"
	"todo-list/models"
)

var rnd *renderer.Render

var posts map[string]*models.Post

/* type DataForm struct {
    Name string `json:"username"`
    Email string `json:"email"`
    Content string `json:"content"`
} */

func init() {
	opts := renderer.Options {
		ParseGlobPattern: "./public/*.html",
	}
	rnd = renderer.New(opts)
}

func index(w http.ResponseWriter, r *http.Request) {
    data := struct {
            Posts map[string]*models.Post
        } {Posts: posts}
	rnd.HTML(w, http.StatusOK, "index", data)
}

func addtask(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "addtask", nil)
}

func userData(w http.ResponseWriter, r *http.Request) {
    username := r.PostFormValue("username")
    email := r.PostFormValue("email")
    content := r.PostFormValue("content")

    if (username == "" || email == "" || content == "") {
        http.Redirect(w, r, "/", 301)
    }

    post := models.NewPost(username, email, content)
    posts[post.Username] = post

    http.Redirect(w, r, "/addtask", 301)

    /* dataForm := DataForm{username, email, content}
    jsonData, err := json.Marshal(dataForm)
    if (err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData) */
}

func main() {
	mux := http.NewServeMux()
	posts = make(map[string]*models.Post, 0)
	fmt.Println(posts)
    mux.HandleFunc("/", index)
    mux.HandleFunc("/addtask", addtask)
    mux.HandleFunc("/userData", userData)
    port := ":8080"
    fmt.Println("Listening on port ", port)
    http.ListenAndServe(port, mux)
}
