package api

import (
	"encoding/json"
	"fmt"
	"github.com/didiercrunch/doormanserver/doormen"
	"github.com/didiercrunch/doormanserver/shared"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func CreateDoorman(w http.ResponseWriter, request *http.Request) {
	wdef := new(doormen.DoormanDefinition)
	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(wdef); err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "cannot create doorman definition \n", err)
		return
	}
	wdef.Id = bson.NewObjectId()

	if err := wdef.Validate(GetUser(request)); err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "New doorman is invalid \n", err)
		return
	}
	if err := CreateDoormanInDatabase(wdef); err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Server error \n", err)
		return
	}
	go publisher.Emit(wdef.Id, wdef)
	w.Header().Set("location", "/api/doormen/"+shared.ObjectIdToPublicId(wdef.Id))
	w.WriteHeader(201)
}
