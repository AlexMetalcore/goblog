package main

import (
    "os"
	"fmt"
	"net/http"
	"strconv"
	"database/sql"
	"github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/thedevsaddam/renderer"
    _ "github.com/go-sql-driver/mysql"
    "github.com/AndyEverLie/go-pagination-bootstrap"
)

var rnd *renderer.Render
var database *sql.DB
var dbName string

type Post struct {
    Id  string
    Username  string
    Email string
    Content string
}

func init() {
	opts := renderer.Options {
		ParseGlobPattern: "./public/*.html",
	}
	rnd = renderer.New(opts)
}

func index(w http.ResponseWriter, r *http.Request) {

    var count int
    countPosts, err := database.Prepare("SELECT COUNT(*) as count FROM " + dbName + ".posts")

    if (err != nil) {
       fmt.Println(err)
    }

    err = countPosts.QueryRow().Scan(&count)

    if (err != nil) {
        fmt.Println(err)
    }

    /* page := "1"
    if (r.FormValue("page") != "") {
        page = r.FormValue("page")
    }

    pagePrepare, err := strconv.Atoi(page)

    if (err != nil) {
        fmt.Println(err)
    }

    limit := 3
    offset := limit * (pagePrepare - 1)

    rows, err := database.Query("SELECT * FROM " + dbName + ".posts ORDER BY id DESC LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit) + "") */

    rows, err := database.Query("SELECT * FROM " + dbName + ".posts")

    if (err != nil) {
        fmt.Println(err)
    }

    defer rows.Close()
    postsData := []Post{}

    for rows.Next() {
        post := Post{}
        err := rows.Scan(&post.Id, &post.Username, &post.Email, &post.Content)
        if (err != nil) {
            fmt.Println(err)
            continue
        }
        postsData = append(postsData, post)
    }

    current, err := strconv.Atoi(r.FormValue("page"))
    if (err != nil) {
        fmt.Println(err)
    }

    //pager := pagination.New(count, limit, current, "/")
    pager := pagination.New(count, 1, current, "/")

    data := struct {
        Posts []Post
        Render *pagination.Pagination
    } {Posts: postsData, Render: pager}

	rnd.HTML(w, http.StatusOK, "home", data)
}

func addPost(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "addPost", nil)
}

func editPost(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")
    row := database.QueryRow("SELECT * FROM " + dbName + ".posts WHERE id = ?", id)
    post := Post{}
    err := row.Scan(&post.Id, &post.Username, &post.Email, &post.Content)

    data := struct {
        Post Post
    } {Post: post}

    if (err != nil) {
        fmt.Println(err)
        http.Error(w, http.StatusText(404), http.StatusNotFound)
    } else {
        rnd.HTML(w, http.StatusOK, "editPost", data)
    }
}

func deletePost(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")
    if (id != "") {
        row := database.QueryRow("SELECT * FROM " + dbName + ".posts WHERE id = ?", id)
        post := Post{}
        err := row.Scan(&post.Id, &post.Username, &post.Email, &post.Content)

        if (err != nil) {
           fmt.Println(err)
           http.Error(w, http.StatusText(404), http.StatusNotFound)
        }

        if (post.Id != "") {
            _, err := database.Exec("DELETE FROM " + dbName + ".posts where id = ?", id)
            if (err != nil) {
               http.Error(w, http.StatusText(404), http.StatusNotFound)
            }
            http.Redirect(w, r, "/", 301)
        }
    } else {
        http.NotFound(w, r)
    }
}

func userData(w http.ResponseWriter, r *http.Request) {
    username := r.PostFormValue("username")
    email := r.PostFormValue("email")
    content := r.PostFormValue("content")

    if (username == "" || email == "" || content == "") {
        http.Redirect(w, r, "/addPost", 301)
    } else {
        if (r.PostFormValue("id") != "") {
            id := r.PostFormValue("id")
            row := database.QueryRow("SELECT * FROM " + dbName + ".posts WHERE id = ?", id)
            post := Post{}
            err := row.Scan(&post.Id, &post.Username, &post.Email, &post.Content)

            if (err != nil) {
               fmt.Println(err)
               http.Error(w, http.StatusText(404), http.StatusNotFound)
            }

            if (post.Id != "") {
                _, err = database.Exec("UPDATE " + dbName + ".posts set username=?, email=?, content = ? where id = ?",
                username, email, content, post.Id)
            } else {
                _, err = database.Exec("INSERT INTO " + dbName + ".posts (username, email, content) VALUES (?, ?, ?)",
                username, email, content)
            }
        } else {
            _, err := database.Exec("INSERT INTO " + dbName + ".posts (username, email, content) VALUES (?, ?, ?)",
            username, email, content)

            if (err != nil) {
                fmt.Println(err)
            }
        }

        http.Redirect(w, r, "/", 301)
    }
}

func main() {
    e := godotenv.Load()

	if (e != nil) {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

    db, err := sql.Open("mysql", "" + username + ":" + password + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "")

    if (err != nil) {
        fmt.Println(err)
    }

    database = db
    defer db.Close()

	mux := mux.NewRouter()
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
