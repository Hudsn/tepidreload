package tepidreload

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed script.tmpl
var scriptFS embed.FS

func MakeHandlers(checkPath string, config Config) (http.HandlerFunc, http.HandlerFunc) {
	t, err := template.ParseFS(scriptFS, "*.tmpl")
	if err != nil {
		panic(err)
	}

	data := map[string]any{
		"Path":     checkPath,
		"Interval": config.TickIntervalMS,
	}

	scriptHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		t.ExecuteTemplate(w, "script.tmpl", data)
	}

	boolHandler := func(w http.ResponseWriter, r *http.Request) {
		isChanged, err := CheckFileMods(config)
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
