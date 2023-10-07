package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hudsn/tepidreload"
)

func main() {
	portNum := flag.String("port", "3000", "port number for the reloader to serve your files on")
	watchIntevalMS := flag.Int("interval", 250, "how many MS to wait between checking for updates")
	excludeExt := flag.String("exclude-ext", "", "comma-delimited list of file extensions to exclude")
	excludeDir := flag.String("exclude-dir", "", "comma-delimited list of directory names to exclude")
	excludeFiles := flag.String("exclude-files", "", "comma-delimited list of specific files to exclude")
	watchPath := flag.String("path", "./", "relative path to directory to monitor for reloads")
	flag.Parse()

	// convert exclusion opts into slices so we can pass to configs
	exDirList := strings.Split(*excludeDir, ",")
	trimmedDirList := []string{}
	for _, entry := range exDirList {
		trimmed := strings.TrimSpace(entry)
		trimmedDirList = append(trimmedDirList, trimmed)
	}
	exExtList := strings.Split(*excludeExt, ",")
	trimmedExtList := []string{}
	for _, entry := range exExtList {
		trimmed := strings.TrimSpace(entry)
		trimmedExtList = append(trimmedDirList, trimmed)
	}
	exFileList := strings.Split(*excludeFiles, ",")
	trimmedFileList := []string{}
	for _, entry := range exFileList {
		trimmed := strings.TrimSpace(entry)
		trimmedExtList = append(trimmedFileList, trimmed)
	}

	tepidConfig := tepidreload.NewConfig(
		tepidreload.WithWatchPath(*watchPath),
		tepidreload.WithExcludeDirs(trimmedDirList...),
		tepidreload.WithExcludeExtensions(trimmedExtList...),
		tepidreload.WithExcludeFiles(trimmedFileList...),
		tepidreload.WithInterval(*watchIntevalMS),
	)

	// Init endpoints for serving script and polling FS
	script, checker := tepidreload.MakeHandlers("/tepid", *portNum, tepidConfig)

	r := chi.NewRouter()

	staticHandler, ok := http.StripPrefix("/tepidstatic/", serveStatic(*watchPath)).(http.HandlerFunc)
	if !ok {
		log.Fatal("Unable to initialize static handler.")
	}

	r.Get("/tepid.js", script)
	r.Get("/tepid", checker)
	r.Get("/tepidstatic/*", staticHandler)
	r.Get("/*", tepidreload.IframeHandler())

	fmt.Printf("Starting the reload server on :%s\n", *portNum)
	fmt.Println("Serving files from ", *watchPath)
	http.ListenAndServe(fmt.Sprintf(":%s", *portNum), r)
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
