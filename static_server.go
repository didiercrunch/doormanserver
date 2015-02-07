package main

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

const STATIC_DIR_PATH = "static/"

func serveStatic(router *mux.Router) {
	handler := func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		filepath := "/" + vars["path"]
		w.Header().Set("Cache-Control", "public, max-age=43200")
		http.ServeFile(w, request, path.Join(STATIC_DIR_PATH, filepath))
	}
	router.HandleFunc("/{path:.*}", handler)
}
