package main

import (
	"fmt"
	"github.com/didiercrunch/doormanserver/api"
	"github.com/didiercrunch/doormanserver/shared"
	"log"
	"net/http"
	"path"
	"strconv"
)

const STATIC_DIR_PATH = "static/"

func GetAddressToServe() string {
	params := shared.GetParams()
	return params.Host + ":" + strconv.Itoa(params.Port)
}

func serveStatic2(w http.ResponseWriter, request *http.Request) {
	filepath := request.URL.Path
	w.Header().Set("Cache-Control", "private, max-age=0, no-cache")
	http.ServeFile(w, request, path.Join(STATIC_DIR_PATH, filepath))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/", api.Create("/api"))
	mux.HandleFunc("/", serveStatic2)

	address := GetAddressToServe()
	fmt.Println("server running at", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}
