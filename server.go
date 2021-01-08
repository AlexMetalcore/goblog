package main

import (
	"fmt"
	"net/http"
)

func serverData(w http.ResponseWriter, r *http.Request) {
        var id string = r.FormValue("id")
		fmt.Fprintf(w, "Hello World, " + "Alex \n")
		fmt.Fprintf(w, id + "\n")
		if (id != "") {
		    fmt.Fprint(w, len(id))
		}
}

func main() {
	http.HandleFunc("/", serverData)
	http.ListenAndServe(":8080", nil)
}
