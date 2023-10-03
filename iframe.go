package tepidreload

import (
	"html/template"
	"log"
	"net/http"
)

// TODO:
// load iframe template at path /*
// add logic to iframe template that updates the src to staticdir/urlpath

func IframeHandler() http.HandlerFunc {
	type TemplateData struct {
		StaticPath string
	}

	t, err := template.ParseFiles("iframe.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		data := TemplateData{
			StaticPath: r.URL.Path,
		}

		t.ExecuteTemplate(w, t.Name(), data)
	}
}
