package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didiercrunch/doormanserver/doormen"
	"github.com/didiercrunch/doormanserver/inmemoryconnector"
)

func init() {
	conn = inmemoryconnector.NewInMemoryDBConnector()
}

func TestGetAllDoormen(t *testing.T) {
	conn = inmemoryconnector.NewInMemoryDBConnector()
	conn.Save(doormen.QuickNewDoormanDefinition("foo", "0.25", "0.75"))
	conn.Save(doormen.QuickNewDoormanDefinition("bar", "0.5", "0.5"))

	req, err := http.NewRequest("GET", "http://bigtits.com/api/doormen", nil)
	if err != nil {
		t.Error(err)
		return
	}

	w := httptest.NewRecorder()
	GetAllDoormen(w, req)

	d := json.NewDecoder(w.Body)
	m := make(map[string][]*doormen.DoormanId)

	if err := d.Decode(&m); err != nil {
		t.Error(err)
		return
	}
	if len(m["doormen"]) != 2 {
		t.Error()
	}
}

func TestGetAllDoormenEmpty(t *testing.T) {
	conn = inmemoryconnector.NewInMemoryDBConnector()
	req, err := http.NewRequest("GET", "http://bigtits.com/api/doormen", nil)
	if err != nil {
		t.Error(err)
		return
	}

	w := httptest.NewRecorder()
	GetAllDoormen(w, req)

	if s := w.Body.String(); s != `{"doormen" :[]}` {
		t.Error(s)
	}
}
