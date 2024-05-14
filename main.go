package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var (
	template_cache = NewTemplateCache()
	music_list     map[string]*Markdown
)

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
	if r.URL.Path == "/music/" {
		renderTemplate(w, r, "music/index.tmpl", "index", music_list)
	} else {

		music, ok := music_list[r.URL.Path[len("/music/"):]+".md"]
		if ok {
			renderTemplate(w, r, "music/music.tmpl", "music", music)
		} else {
			http.ServeFile(w, r, "static/"+r.URL.Path[1:])
		}
	}
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/gallery" {
		file, err := LoadFile("gallery/gallery.json")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var images []Image
		json.Unmarshal(file, &images)

		renderTemplate(w, r, "gallery/index.tmpl", "index", images)
	}
}

func initMusicList() error {
	music, err := ReadDir("music")
	if err != nil {
		return err
	}

	music_list = make(map[string]*Markdown)
	for _, entry := range music {
		if !strings.Contains(entry.Name(), ".md") {
			continue
		}

		md, err := ParseMarkdown("music/" + entry.Name())
		if err != nil {
			log.Println(err)
			continue
		}

		music_list[entry.Name()] = md
	}

	return nil
}

func main() {
	err := initMusicList()
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/music/", musicHandler)
	http.HandleFunc("/gallery/", galleryHandler)

	log.Fatal(http.ListenAndServe(":8484", nil))
}
