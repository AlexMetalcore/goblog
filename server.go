package main

import (
	"fmt"
	"net/http"
	"github.com/thedevsaddam/renderer"
	"encoding/json"
)

var rnd *renderer.Render

type DataForm struct {
    Name string `json:"username"`
    Email string `json:"email"`
    Content string `json:"content"`
}

func init() {
	opts := renderer.Options {
		ParseGlobPattern: "./public/*.html",
	}
	rnd = renderer.New(opts)
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
    dataForm := DataForm{username, email, content}
    jsonData, err := json.Marshal(dataForm)
    if (err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}

func main() {
	mux := http.NewServeMux()
    mux.HandleFunc("/", addtask)
    mux.HandleFunc("/userData", userData)
    port := ":8080"
    fmt.Println("Listening on port ", port)
    http.ListenAndServe(port, mux)
}
