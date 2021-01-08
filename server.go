package main

import (
	"fmt"
	"net/http"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options {
		ParseGlobPattern: "./public/*.html",
	}

	rnd = renderer.New(opts)
}

func home(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")
	rnd.HTML(w, http.StatusOK, "home", struct{Data string}{Data: id})
}

func about(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "about", nil)
}

/* func serverData(w http.ResponseWriter, r *http.Request) {
        var id string = r.FormValue("id")
		fmt.Fprintf(w, "Hello World, " + "Alex \n")
		fmt.Fprintf(w, id + "\n")
		if (id != "") {
		    fmt.Fprint(w, len(id))
		}
} */

func main() {
	mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/about", about)
    port := ":8080"
    fmt.Println("Listening on port ", port)
    http.ListenAndServe(port, mux)
}
