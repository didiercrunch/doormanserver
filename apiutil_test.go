package main

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestWrite404Error(t *testing.T) {
	rw := httptest.NewRecorder()
	Write404Error(rw)
	if rw.Code != 404 {
		t.Error()
	}
}

func TestWrite500Error(t *testing.T) {
	rw := httptest.NewRecorder()
	Write500Error(rw, errors.New("blah"))
	if rw.Code != 500 {
		t.Error()
	}
	if string(rw.Body.Bytes()) != "server error" {
		t.Error()
	}
}

func TestWrite400Error(t *testing.T) {
	rw := httptest.NewRecorder()
	Write400Error(rw, errors.New("blah"))
	if rw.Code != 400 {
		t.Error()
	}
	if string(rw.Body.Bytes()) != "blah" {
		t.Error()
	}
}
