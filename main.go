package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var template_cache = NewTemplateCache()

func loadFile(path string) ([]byte, error) {
	file, err := os.ReadFile("static/" + path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return file, nil
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, tmpl_name string, data any) {
	t, err := template_cache.Get(tmpl)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}

	err = t.ExecuteTemplate(w, tmpl_name, data)
	if err != nil {
		log.Println(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		renderTemplate(w, r, "index.tmpl", "page", nil)
		return
	}

	idx := strings.Index(r.URL.Path, ".")
	if idx == -1 {
		renderTemplate(w, r, r.URL.Path[1:]+".tmpl", "page", nil)
	} else {
		http.ServeFile(w, r, "static/"+r.URL.Path[1:])
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/blog" {
		renderTemplate(w, r, "blog/index.tmpl", "index", nil)
		return
	}
}

func musicHandler(w http.ResponseWriter, r *http.Request) {
	file, err := loadFile("music/music.json")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var music map[string]interface{}
	err = json.Unmarshal(file, &music)
	if err != nil {
		log.Println(err)
	}

	if r.URL.Path == "/music/" {
		renderTemplate(w, r, "music/index.tmpl", "index", music)
	} else {
		renderTemplate(w, r, r.URL.Path[1:], "page", music)
	}
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/gallery" {
		file, err := loadFile("gallery/gallery.json")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var images []Image
		json.Unmarshal(file, &images)

		renderTemplate(w, r, "gallery/index.tmpl", "index", images)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/music/", musicHandler)
	http.HandleFunc("/gallery/", galleryHandler)

	log.Fatal(http.ListenAndServe(":8484", nil))
}
