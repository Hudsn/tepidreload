package main

import (
	"bytes"
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
	templates    *template.Template
	devTemplates tepidreload.DevTemplates
}

// An example of how you might use the package in a Go project
func main() {
	app := MyApp{}
	lport := flag.String("port", "9001", "listen port for app")
	flag.BoolVar(&app.Config.devTemplates.IsDev, "dev", true, "whether it's a dev environment")
	flag.Parse()

	// generate templates based on env mode
	// dev generates from files on the system to enable hot reload (we want to watch the tmpl files we're editing to trigger them; if we use embedded tmpl files it won't trigger reload because we're not editing those.)
	switch app.Config.devTemplates.IsDev {
	case true:
		app.Config.devTemplates.MakeLocalDevTemplates("example/static", "tmpl")
	case false:
		tt, err := template.ParseFS(exampleFS, "static/*.tmpl")
		if err != nil {
			log.Fatalf("Unable to parse templates: %v", err)
		}
		app.Config.templates = tt
	}

	// create a new config specifying the root of our path we want to watch
	// note: you can instead specify an embedded FS with 'tepidreload.WithEmbedFS()' instead of WithWatchPath, just don't use both as args to NewConfig
	watchPath := "example"
	tepidConfig := tepidreload.NewConfig(tepidreload.WithWatchPath(watchPath))

	// specify the endpoint that we'll use to check for reloads. We only need this value because it'll be dynamically added to the polling script.
	// we get back two handler funcs to pass to any router that's compatible with net/http package
	handlerEndpointPath := "/tepid"
	scriptHandler, checkHandler := tepidreload.MakeHandlers(handlerEndpointPath, *lport, tepidConfig)

	app.Router = chi.NewRouter()
	if app.Config.devTemplates.IsDev {
		// this first path can be anything you want, just make sure you reference it when you add the script tag to your HTML templates (see static/example.tmpl)
		app.Router.Get("/tepid.js", scriptHandler)
		// this second handler should match the arg you passed to MakeHandlers
		app.Router.Get(handlerEndpointPath, checkHandler)
	}

	// normal routes
	// NOTE: we conditionally add the polling script and endpoint by using app.Config to read "IsDev" and pass that to templates so we only have hot reload in dev (see: defaultTemplateData and how it's used in exampleHandler())
	// Try using "go run examples/main.go -dev=false" and then check to see that the script is no longer rendered, and that /tepid returns a 404 error.
	app.Router.Get("/", app.exampleHandler())

	log.Println("Listening on :9001")
	fmt.Println()
	fmt.Println("Navigate to localhost:9001 to see the Hello World page...")
	fmt.Println("then make some changes to example.tmpl and the page should reload!")
	http.ListenAndServe(":9001", app.Router)
}

type templateData struct {
	IsDev bool
}

func (a *MyApp) defaultTemplateData() templateData {
	return templateData{
		IsDev: a.Config.devTemplates.IsDev,
	}
}

//go:embed static/*
var exampleFS embed.FS

func (a *MyApp) exampleHandler() http.HandlerFunc {
	data := a.defaultTemplateData()

	return func(w http.ResponseWriter, r *http.Request) {
		a.render(w, "example.tmpl", data)
	}
}

func (a *MyApp) render(w http.ResponseWriter, templateName string, data any) {
	wbuf := bytes.NewBuffer([]byte{})

	// render template directly from local file if it's a dev environment
	if a.Config.devTemplates.IsDev {

		templatePath, found := a.Config.devTemplates.GetLocalTemplate(templateName)
		if !found {
			fmt.Println("Templates not found")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles(templatePath)
		if err != nil {
			fmt.Println("Failed to parse template")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = t.ExecuteTemplate(wbuf, templateName, data)
		if err != nil {
			fmt.Println("Failed to execute template")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		wbuf.WriteTo(w)
		return
	}

	// normal rendering from template cache:
	err := a.Config.templates.ExecuteTemplate(wbuf, templateName, data)
	if err != nil {
		fmt.Println("Failed to execute template")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	wbuf.WriteTo(w)
	return

}
