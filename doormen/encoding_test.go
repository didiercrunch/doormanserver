package doormen

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
)

func TestJsonEncoding(t *testing.T) {
	dmd := QuickNewDoormanDefinition("test json", "1/4", "3/4")
	data, err := json.Marshal(dmd)
	if err != nil {
		t.Error(err)
	}
	dmd2 := new(DoormanDefinition)
	if err = json.Unmarshal(data, dmd2); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(dmd2, dmd) {
		t.Error()
	}
}

func TestBsonEncoding(t *testing.T) {
	dmd := QuickNewDoormanDefinition("test bson", "1/4", "3/4")
	data, err := bson.Marshal(dmd)
	if err != nil {
		t.Error(err)
	}
	dmd2 := new(DoormanDefinition)
	if err = bson.Unmarshal(data, dmd2); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(dmd2, dmd) {
		t.Error()
	}
}
