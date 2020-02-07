package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gosidekick/bin2go/example/assets"
)

var mimeType map[string]string

func showPage(w http.ResponseWriter, name string, data interface{}) {
	f, _ := assets.GetBytes(name)
	t, err := template.New(name).Parse(string(f))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Title   string
		Comment string
	}
	d := data{
		Title:   "Golang",
		Comment: "The Gophers programming language ;D",
	}

	showPage(w, "index.html", d)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	f := mux.Vars(r)
	name := f["filename"]
	b, err := assets.GetBytes(name)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404 file not found"))
		return
	}

	l := len(b)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	t, ok := mimeType[filepath.Ext(name)]
	if !ok {
		if l > 512 {
			l = 512
		}
		t = http.DetectContentType(b[:l])

	}
	w.Header().Set("Content-Type", t)
	w.Write(b)
}

func main() {
	mimeType = make(map[string]string)
	mimeType[".css"] = "text/css"
	mimeType[".svg"] = "image/svg+xml"
	mimeType[".png"] = "image/png"

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", handleMain).Methods(http.MethodGet)
	r.HandleFunc("/{filename}", handleFile).Methods(http.MethodGet)

	fmt.Println("access the page at localhost:8080")
	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
