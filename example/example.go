package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hudsn/tepidreload"
)

type MyApp struct {
	Router *chi.Mux
	Config config
}

type config struct {
	isDev     bool
	templates *template.Template
}

// An example of how you might use the package in a Go project
func main() {
	app := MyApp{}
	lport := flag.String("port", "9001", "listen port for app")
	flag.BoolVar(&app.Config.isDev, "dev", true, "whether it's a dev environment")
	flag.Parse()

	tt, err := template.ParseFS(exampleFS, "*.tmpl")
	if err != nil {
		log.Fatalf("Unable to parse templates: %v", err)
	}

	app.Config.templates = tt

	// create a new config specifying the root of our path we want to watch
	// note: you can instead specify an embedded FS with 'tepidreload.WithEmbedFS()' instead of WithWatchPath, just don't use both as args to NewConfig
	watchPath := "test_static"
	tepidConfig := tepidreload.NewConfig(tepidreload.WithWatchPath(watchPath))

	// specify the endpoint that we'll use to check for reloads. We only need this value because it'll be dynamically added to the polling script.
	// we get back two handler funcs to pass to any router that's compatible with net/http package
	scriptHandler, checkHandler := tepidreload.MakeHandlers("/tepid", *lport, tepidConfig)

	app.Router = chi.NewRouter()

	if app.Config.isDev {
		// this first path can be anything you want, just make sure you reference it when you add the script tag to your HTML templates (see example.tmpl)
		app.Router.Get("/tepid.js", scriptHandler)
		// this second handler should match the arg you passed to MakeHandlers
		app.Router.Get("/tepid", checkHandler)
	}

	// normal routes
	// NOTE: we conditionally add the polling script and endpoint by using app.Config to read "isDev" and pass that to templates so we only have hot reload in dev (see: defaultTemplateData and how it's used in exampleHandler())
	// Try using "go run examples/main.go -dev=false" and then check to see that the script is no longer rendered, and that /tepid returns a 404 error.
	app.Router.Get("/", app.exampleHandler())

	log.Println("Listening on :9001")
	fmt.Println()
	fmt.Println("Navigate to localhost:9001 to see the Hello World page...")
	fmt.Println("then make some changes to files in the test_static directory and the page should reload!")
	http.ListenAndServe(":9001", app.Router)
}

type templateData struct {
	IsDev bool
}

func (a *MyApp) defaultTemplateData() templateData {
	return templateData{
		IsDev: a.Config.isDev,
	}
}

//go:embed example.tmpl
var exampleFS embed.FS

func (a *MyApp) exampleHandler() http.HandlerFunc {
	data := a.defaultTemplateData()

	return func(w http.ResponseWriter, r *http.Request) {
		err := a.Config.templates.ExecuteTemplate(w, "example.tmpl", data)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
