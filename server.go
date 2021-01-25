package main

import (
	"fmt"
	"net/http"
	"github.com/thedevsaddam/renderer"
	"github.com/gorilla/mux"
	"crypto/rand"
	"os"
    "github.com/joho/godotenv"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var rnd *renderer.Render
var database *sql.DB
var dbName string
var posts map[string]*Post

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
    rows, err := database.Query("select * from " + dbName + ".posts")

    if err != nil {
        fmt.Println(err)
    }

    defer rows.Close()
    postsData := []Post{}

    for rows.Next(){
        post := Post{}
        err := rows.Scan(&post.Id, &post.Username, &post.Email, &post.Content)
        if err != nil{
            fmt.Println(err)
            continue
        }
        postsData = append(postsData, post)
    }

    data := struct {
        Posts []Post
    } {Posts: postsData}

	rnd.HTML(w, http.StatusOK, "home", data)
}

func addPost(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "addPost", nil)
}

func editPost(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")
    row := database.QueryRow("select * from "+dbName+".posts WHERE id = ?", id)
    post := Post{}
    err := row.Scan(&post.Id, &post.Username, &post.Email, &post.Content)

    data := struct {
        Post Post
    } {Post: post}

    if err != nil{
        fmt.Println(err)
        http.Error(w, http.StatusText(404), http.StatusNotFound)
    } else {
        rnd.HTML(w, http.StatusOK, "editPost", data)
    }
}

func deletePost(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")
    if (id != "") {
        row := database.QueryRow("select * from " + dbName + ".posts WHERE id = ?", id)
        post := Post{}
        err := row.Scan(&post.Id, &post.Username, &post.Email, &post.Content)
        if (err != nil) {
           fmt.Println(err)
           http.Error(w, http.StatusText(404), http.StatusNotFound)
        }
        fmt.Println(post.Id)
        if (post.Id != "") {
            _, err := database.Exec("delete from " + dbName + ".posts where id = ?", id)
            if (err != nil) {
               fmt.Println(err)
               http.Error(w, http.StatusText(404), http.StatusNotFound)
            }
            http.Redirect(w, r, "/", 301)
        } else {
           http.NotFound(w, r)
        }
    }
    http.NotFound(w, r)
}

func userData(w http.ResponseWriter, r *http.Request) {
    username := r.PostFormValue("username")
    email := r.PostFormValue("email")
    content := r.PostFormValue("content")

    if (username == "" || email == "" || content == "") {
        http.Redirect(w, r, "/", 301)
    } else {
        if (r.PostFormValue("id") != "") {
            id := r.PostFormValue("id")
            row := database.QueryRow("select * from " + dbName + ".posts WHERE id = ?", id)
            post := Post{}
            err := row.Scan(&post.Id, &post.Username, &post.Email, &post.Content)

            if (err != nil) {
               fmt.Println(err)
               http.Error(w, http.StatusText(404), http.StatusNotFound)
            }

            if (post.Id != "") {
                _, err = database.Exec("update " + dbName + ".posts set username=?, email=?, content = ? where id = ?", username, email, content, post.Id)
            } else {
                _, err = database.Exec("insert into " + dbName + ".posts (username, email, content) values (?, ?, ?)", username, email, content)
            }
        } else {
            _, err := database.Exec("insert into " + dbName + ".posts (username, email, content) values (?, ?, ?)", username, email, content)

            if (err != nil) {
                fmt.Println(err)
            }
        }

        http.Redirect(w, r, "/addPost", 301)
    }
}

func main() {
    e := godotenv.Load()

	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

    db, err := sql.Open("mysql", ""+username+":"+password+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"")

    if err != nil {
        fmt.Println(err)
    }
    database = db
    defer db.Close()

	mux := mux.NewRouter()
	posts = make(map[string]*Post, 0)
	fmt.Println(posts)
	router := mux.StrictSlash(true)
    router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
    mux.HandleFunc("/", index)
    mux.HandleFunc("/addPost", addPost)
    mux.HandleFunc("/editPost", editPost)
    mux.HandleFunc("/deletePost", deletePost)
    mux.HandleFunc("/userData", userData)
    port := ":8080"
    fmt.Println("Listening on port ", port)
    http.ListenAndServe(port, mux)
}
