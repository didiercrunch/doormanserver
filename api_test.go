package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didiercrunch/doormanserver/inmemoryconnector"
	"github.com/didiercrunch/doormanserver/doormen"
	"github.com/gorilla/mux"
)

func init() {
	conn = inmemoryconnector.NewInMemoryDBConnector()
}

func OkHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "ok")
}

func assertBodyIs(t *testing.T, res *http.Response, expt string) {
	if body, err := ioutil.ReadAll(res.Body); err != nil {
		t.Error(err)
	} else if s := string(body); s != expt {
		t.Error("expected ", expt, " but received ", s)
	}
}

func TestPOSTToDoormen(t *testing.T) {
	api := &Api{CreateDoorman: OkHandler}
	r := mux.NewRouter()
	api.ServeApi(r)

	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Post(ts.URL+"/doormen", "", bytes.NewBufferString("boby"))
	if err != nil {
		t.Error(err)
		return
	}
	assertBodyIs(t, res, "ok")
}

func TestGETToDoormen(t *testing.T) {
	api := &Api{GetAllDoormen: OkHandler}

	r := mux.NewRouter()
	api.ServeApi(r)

	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/doormen")
	if err != nil {
		t.Error(err)
		return
	}
	assertBodyIs(t, res, "ok")
}

func TestGetAllDoormenEmpty(t *testing.T) {
	conn = inmemoryconnector.NewInMemoryDBConnector()
	api := &Api{GetAllDoormen: GetAllDoormen}

	r := mux.NewRouter()
	api.ServeApi(r)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/doormen")
	if err != nil {
		t.Error(err)
		return
	}

	if body, err := ioutil.ReadAll(res.Body); err != nil {
		t.Error(err)
	} else if s := string(body); s != `{"doormen" :[]}` {
		t.Error(s)
	}
}

func TestGetAllDoormen(t *testing.T) {
	conn = inmemoryconnector.NewInMemoryDBConnector()
	conn.Save(doormen.QuickNewDoormanDefinition("foo", 0.25, 0.75))
	conn.Save(doormen.QuickNewDoormanDefinition("bar", 0.5, 0.5))

	api := &Api{GetAllDoormen: GetAllDoormen}
	r := mux.NewRouter()
	api.ServeApi(r)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/doormen")
	if err != nil {
		t.Error(err)
		return
	}
	d := json.NewDecoder(res.Body)
	m := make(map[string][]*doormen.DoormanDefinition)

	if err := d.Decode(&m); err != nil {
		t.Error(err)
		return
	}
	if len(m["doormen"]) != 2 {
		t.Error()
	}

}
