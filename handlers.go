package tepidreload

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed iframe.tmpl
var iframeFS embed.FS

//go:embed script.tmpl
var scriptFS embed.FS

// Returns two handlers: One to serve a generated script file to poll for reloads, and the second is the endpoint that the polling script checks against.
func MakeHandlers(checkPath string, config Config) (http.HandlerFunc, http.HandlerFunc) {
	t, err := template.ParseFS(scriptFS, "*.tmpl")
	if err != nil {
		panic(err)
	}

	type Data struct {
		Path     string
		Interval int
	}

	data := Data{
		Path:     checkPath,
		Interval: config.TickIntervalMS,
	}

	scriptHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		t.ExecuteTemplate(w, "script.tmpl", data)
	}

	boolHandler := func(w http.ResponseWriter, r *http.Request) {
		isChanged, err := checkFileMods(config)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Expires", "0")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		if isChanged == true {
			fmt.Fprint(w, "true")
			return
		}

		fmt.Fprint(w, "false")
	}

	return scriptHandler, boolHandler
}

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
