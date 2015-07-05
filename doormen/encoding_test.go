package doormen

import (
	"encoding/base64"
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

func TestJsonEncodingWithEmptyIds(t *testing.T) {
	dmd := QuickNewDoormanDefinition("test json", "1/4", "3/4")
	dmd.Id = ""
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

func TestJsonIdFormatToBeCompatibleWithSiphashKey(t *testing.T) {
	dmd := QuickNewDoormanDefinition("test json", "1/4", "3/4")
	data, err := json.Marshal(dmd)
	if err != nil {
		t.Error(err)
		return
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(data, &m); err != nil {
		t.Error(err)
		return
	}
	if b, err := base64.URLEncoding.DecodeString(m["id"].(string)); err != nil {
		t.Error(err)
		return
	} else if len(b) != 16 {
		t.Error("siphash requires a 16 bytes (128 bits) key not ", len(b))
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

func TestDoormanIdJsonEncoding(t *testing.T) {
	did := &DoormanId{bson.NewObjectId(), "hahaha"}
	data, err := json.Marshal(did)
	if err != nil {
		t.Error(err)
	}
	did2 := new(DoormanId)
	if err = json.Unmarshal(data, did2); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(did2, did) {
		t.Error()
	}

}
