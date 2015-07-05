package api

import (
	"bytes"
	"github.com/didiercrunch/doormanserver/inmemoryconnector"
	"github.com/didiercrunch/doormanserver/shared"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func GetIdFromLocation(p string) bson.ObjectId {
	if id, err := shared.PublicIdToObjectId(path.Base(p)); err != nil {
		panic(err)
	} else {
		return id
	}

}

func TestCreateDoorman(t *testing.T) {
	conn = inmemoryconnector.NewInMemoryDBConnector()
	payload := `
	{"name":"name",
	"values":[{"name":"red","probability":"1/5"},{"name":"blue","probability": "4/5"}],
	"emails": ["alice"]
	}`
	req, err := http.NewRequest(
		"POST",
		"http://bigtits.com/api/doormen",
		bytes.NewBufferString(payload),
	)
	if err != nil {
		t.Error(err)
		return
	}
	SetUser(req, "alice")
	w := httptest.NewRecorder()
	CreateDoorman(w, req)
	if w.Code != 201 {
		t.Error("bad status code", w.Code, "\n", w.Body)
	} else if location := w.Header().Get("location"); location == "" {
		t.Error("missing location header", location)
	} else if dor, _ := conn.GetDoorman(GetIdFromLocation(location)); dor == nil {
		t.Error("dorman not created")
	}
}

func TestCreateDoormanWithBadInput(t *testing.T) {
	conn = inmemoryconnector.NewInMemoryDBConnector()
	payload := `{"name":"name","values":[{"name":"red","probability":0.2},{"name":"blue","probability":0.2}]}`
	req, err := http.NewRequest(
		"POST",
		"http://bigtits.com/api/doormen",
		bytes.NewBufferString(payload),
	)
	if err != nil {
		t.Error(err)
		return
	}
	w := httptest.NewRecorder()
	CreateDoorman(w, req)
	if w.Code != 400 {
		t.Error("bad status code", w.Code)
	}
}
