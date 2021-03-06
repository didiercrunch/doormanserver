package api

import (
	"fmt"
	"github.com/didiercrunch/doormanserver/doormen"
	"github.com/didiercrunch/doormanserver/shared"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

func getDoormenIdsAsJson(id *doormen.DoormanId) string {
	publicId := shared.ObjectIdToPublicId(id.Id)
	return fmt.Sprintf(`{"name": "%v", "id": "%v", "url": "/api/doormen/%v"}`, id.Name, publicId, publicId)
}

func getDoormanIdFromRequest(w http.ResponseWriter, request *http.Request) bson.ObjectId {
	id := mux.Vars(request)["id"]
	if id, err := shared.PublicIdToObjectId(id); err != nil {
		return ""
	} else {
		return id
	}
}

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
