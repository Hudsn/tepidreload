package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hudsn/tepidreload"
)

func main() {

	tepidConfig := tepidreload.NewConfig(tepidreload.WithWatchPath("test_static"))
	script, checker := tepidreload.MakeHandlers("/tepid", tepidConfig)

	r := chi.NewRouter()
	r.Get("/tepid.js", script)
	r.Get("/tepid", checker)
	r.Get("/", normal)

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", r)
}

func normal(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/html")
	fmt.Fprint(w, `
	<h1>Hello world</h1>
	<script src="/tepid.js"></script>
	`)
}
