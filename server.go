package main

import (
	"fmt"
	"net/http"
	"github.com/thedevsaddam/renderer"
	"github.com/gorilla/mux"
	"crypto/rand"
	//"encoding/json"
)

var rnd *renderer.Render

var posts map[string]*Post

/* type DataForm struct {
    Name string `json:"username"`
    Email string `json:"email"`
    Content string `json:"content"`
} */

type Post struct {
    Id  string
    Username  string
    Email string
    Content string
}

func NewPost(id, username, email, content string) *Post {
    return &Post{id, username, email, content}
}

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func init() {
	opts := renderer.Options {
		ParseGlobPattern: "./public/*.html",
	}
	rnd = renderer.New(opts)
}

func index(w http.ResponseWriter, r *http.Request) {
    data := struct {
            Posts map[string]*Post
        } {Posts: posts}
	rnd.HTML(w, http.StatusOK, "home", data)
}

func addtask(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "addtask", nil)
}

func userData(w http.ResponseWriter, r *http.Request) {
    id := GenerateId()
    username := r.PostFormValue("username")
    email := r.PostFormValue("email")
    content := r.PostFormValue("content")

    if (username == "" || email == "" || content == "") {
        http.Redirect(w, r, "/", 301)
    }

    post := NewPost(id, username, email, content)
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
	mux := mux.NewRouter()
	posts = make(map[string]*Post, 0)
	fmt.Println(posts)
	router := mux.StrictSlash(true)
    router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
    mux.HandleFunc("/", index)
    mux.HandleFunc("/addtask", addtask)
    mux.HandleFunc("/userData", userData)
    port := ":8080"
    fmt.Println("Listening on port ", port)
    http.ListenAndServe(port, mux)
}
