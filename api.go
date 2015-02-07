package main

import (
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/didiercrunch/doormanserver/doormen"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func getDoormenIdsAsJson(id *doormen.DoormanId) string {
	return fmt.Sprintf(`{"name": "%v", "id": "%v", "url": "/api/doormen/%v"}`, id.Name, id.Id.Hex(), id.Id.Hex())
}

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
	if err := wdef.Validate(); err != nil {
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
	w.Header().Set("location", "/api/doormen/"+wdef.Id.Hex())
	w.WriteHeader(201)
}

func GetAllDoormen(w http.ResponseWriter, request *http.Request) {
	fmt.Fprint(w, `{"doormen" :[`)
	doormenIds := conn.GetAllDoormen()
	buff := <-doormenIds
	if buff == nil {
		fmt.Fprint(w, "]}")
		return
	}
	for doormanId := range doormenIds {
		fmt.Fprint(w, getDoormenIdsAsJson(buff)+`,`)
		buff = doormanId
	}
	fmt.Fprint(w, getDoormenIdsAsJson(buff), "]}")
}

func getDoormanIdFromRequest(w http.ResponseWriter, request *http.Request) bson.ObjectId {
	id := mux.Vars(request)["id"]
	if !bson.IsObjectIdHex(id) {
		return ""
	} else {
		return bson.ObjectIdHex(id)
	}
}

func GetDoorman(w http.ResponseWriter, r *http.Request) {
	id := getDoormanIdFromRequest(w, r)
	if id == "" {
		Write404Error(w)
		return
	}
	doorman, err := conn.GetDoorman(id)
	if err != nil {
		Write500Error(w, err)
	} else if doorman == nil {
		Write404Error(w)
	} else {
		encoder := json.NewEncoder(w)
		encoder.Encode(doorman)
	}
}

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

type Api struct {
	GetAllDoormen http.HandlerFunc
	CreateDoorman  http.HandlerFunc
	GetDoorman     http.HandlerFunc
	UpdateDoorman  http.HandlerFunc
}

func (a *Api) ServeApi(router *mux.Router) {
	router.Path("/doormen").Methods("GET").HandlerFunc(a.GetAllDoormen)
	router.Path("/doormen").Methods("POST").HandlerFunc(a.CreateDoorman)
	router.Path("/doormen/{id}").Methods("GET").HandlerFunc(a.GetDoorman)
	router.Path("/doormen/{id}").Methods("PUT").HandlerFunc(a.UpdateDoorman)
}

func serveApi(router *mux.Router) {
	api := new(Api)
	api.CreateDoorman = CreateDoorman
	api.GetAllDoormen = GetAllDoormen
	api.GetDoorman = GetDoorman
	api.UpdateDoorman = UpdateDoorman
	api.ServeApi(router)
}
