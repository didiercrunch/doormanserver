package main

import (
	"testing"

	"github.com/didiercrunch/doormanserver/inmemoryconnector"
	"github.com/didiercrunch/doormanserver/doormen"
)

func TestCreateDoorman(t *testing.T) {
	wab := doormen.NewDoormanDefinition("some_name")
	conn = inmemoryconnector.NewInMemoryDBConnector()

	if err := CreateDoormanInDatabase(wab); err != nil {
		t.Error(err)
	}

	if err := CreateDoormanInDatabase(wab); err == nil {
		t.Error("should have an error here")
	}
}
