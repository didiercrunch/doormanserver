package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthentication(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if user := GetUser(req); user != "natasha" {
			t.Error(user)
		}
	}

	req, err := http.NewRequest("GET", "http://example.com/foo?user=natasha", nil)
	if err != nil {
		t.Error(err)
		return
	}
	w := httptest.NewRecorder()
	AuthenticationMiddleware(handler)(w, req)
	if user := GetUser(req); user != "" {
		t.Error(user)
	}

}
