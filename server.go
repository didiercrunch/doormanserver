package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/didiercrunch/doormanserver/shared"
	"github.com/gorilla/mux"
)

func createMuxRouter() http.Handler {
	r := mux.NewRouter()
	serveApi(r.PathPrefix("/api").Subrouter())
	serveStatic(r.PathPrefix("/").Subrouter())
	return r
}

func GetAddressToServe() string {
	params := shared.GetParams()
	return params.Host + ":" + strconv.Itoa(params.Port)
}

func main() {
	r := createMuxRouter()
	http.Handle("/", r)
	address := GetAddressToServe()
	fmt.Println("server running at", address)
	if err := http.ListenAndServe(address, createMuxRouter()); err != nil {
		log.Fatal(err)
	}
}
