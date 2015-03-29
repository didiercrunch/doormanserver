package api

import (
	"encoding/json"
	"github.com/didiercrunch/doormanserver/doormen"
	"net/http"
)

func UpdateDoorman(w http.ResponseWriter, request *http.Request) {
	id := getDoormanIdFromRequest(w, request)
	if id == "" || !conn.ExistsId(id) {
		Write404Error(w)
		return
	}

	wdef := new(doormen.DoormanDefinition)
	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(wdef); err != nil {
		Write400Error(w, err)
		return
	}
	wdef.Id = id

	if err := wdef.Validate(); err != nil {
		Write400Error(w, err)
		return
	}

	originalDoorman, err := conn.GetDoorman(id)
	if err != nil {
		Write500Error(w, err)
		return
	}
	if err := originalDoorman.CanBeUpdatedBy(wdef); err != nil {
		Write400Error(w, err)
		return
	}

	if err := UpdateDoormanInDatabase(wdef); err != nil {
		Write500Error(w, err)
		return
	}
	go publisher.Emit(wdef.Id, wdef)
	w.WriteHeader(200)

}
