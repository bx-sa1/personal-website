package personalwebsite

import (
	"log"
	"net/http"
	"os"
	"strings"
)

var template_cache = NewTemplateCache()

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string) {
	t, err := template_cache.Get(tmpl)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	t.Execute(w, nil)
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	idx := strings.Index(r.URL.Path, ".")
	if idx == -1 {
		renderTemplate(w, r, r.URL.Path[1:]+".tmpl")
	} else {
		file, err := os.ReadFile("static/" + r.URL.Path[1:])
		if err != nil {
			http.NotFound(w, r)
		}
		w.Write(file)
	}
}

func main() {
	http.HandleFunc("/*", htmlHandler)

	log.Fatal(http.ListenAndServe(":8484", nil))
}
