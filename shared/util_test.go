package shared

import (
	"encoding/base64"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestObjectIdToString(t *testing.T) {
	b := bson.ObjectIdHex("5596d9b7f6fb5e195e000004")
	s := ObjectIdToPublicId(b)
	key, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		t.Error(err)
		return
	}
	if len(key) != 16 {
		t.Error("key must be of length 16 bytes to be compatible with siphash.", len(key))
	}
}

func TestStringToObjectId(t *testing.T) {
	s := "O725hlWW2bf2-14ZXgAABA=="
	oid, err := PublicIdToObjectId(s)
	if err != nil {
		t.Error(err)
		return
	} else if oid != bson.ObjectIdHex("5596d9b7f6fb5e195e000004") {
		t.Error(oid)
	}
}

func TestValidateId(t *testing.T) {
	if id, err := base64.URLEncoding.DecodeString("cGOZFFWW2bf2-14ZXgAABA=="); err != nil {
		t.Error(err)
	} else if err = validateId(id); err == nil {
		t.Error()
	}

	if id, err := base64.URLEncoding.DecodeString("O725hlWW2bf2-14ZXgAA"); err != nil {
		t.Error(err)
	} else if err = validateId(id); err == nil {
		t.Error()
	}
}
