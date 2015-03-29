package api

import (
	"encoding/json"
	"github.com/didiercrunch/doormanserver/serverstat"
	"github.com/didiercrunch/doormanserver/shared"
	"log"
	"mime"
	"net/http"
)

func GetServerSpecification(w http.ResponseWriter, request *http.Request) {
	s := serverstat.Get(shared.GetParams())
	w.Header().Add("Content-Type", mime.TypeByExtension(".json"))
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(s); err != nil {
		log.Println(err)
	}
}
