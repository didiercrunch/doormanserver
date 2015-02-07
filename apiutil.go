package main

import (
	"fmt"
	"log"
	"net/http"
)

func Write500Error(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprint(w, "server error")
	log.Println(err)
}

func Write404Error(w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprint(w, "url not found")
}

func Write400Error(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	fmt.Fprint(w, err.Error())
}
