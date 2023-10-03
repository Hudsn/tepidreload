package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hudsn/tepidreload"
)

func main() {

	watchPath := flag.String("path", "./", "relative path to directory to monitor for reloads")
	flag.Parse()

	tepidConfig := tepidreload.NewConfig(tepidreload.WithWatchPath(*watchPath))

	script, checker := tepidreload.MakeHandlers("/tepid", tepidConfig)

	r := chi.NewRouter()

	staticHandler, ok := http.StripPrefix("/tepidstatic/", serveStatic(*watchPath)).(http.HandlerFunc)
	if !ok {
		log.Fatal("Unable to initialize static handler.")
	}

	r.Get("/tepid.js", script)
	r.Get("/tepid", checker)
	r.Get("/tepidstatic/*", staticHandler)
	r.Get("/*", tepidreload.IframeHandler())

	fmt.Println("Starting the server on :3000")
	fmt.Println("Serving files from ", *watchPath)
	http.ListenAndServe(":3000", r)
}

func serveStatic(filepath string) http.HandlerFunc {

	dir := http.Dir(filepath)

	return func(w http.ResponseWriter, r *http.Request) {
		// perma-nuke cache so we don't get weird interactions if the user changes the watched directory between uses.
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Cache-Control", "must-revalidate")
		w.Header().Set("Expires", "0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Del("Last-Modified")
		http.FileServer(dir).ServeHTTP(w, r)
	}

}
