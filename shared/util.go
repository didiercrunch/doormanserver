package shared

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"github.com/dchest/siphash"
	"gopkg.in/mgo.v2/bson"
	"os"
)

var idHashKey []byte = []byte("semi-private-key")

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func validateId(id []byte) error {
	if len(id) != 16 {
		return errors.New("valid ids are 16 bytes (126 bit) long.")
	}
	h := siphash.New(idHashKey)
	h.Write(id[4:16])
	if subtle.ConstantTimeCompare(h.Sum(nil)[0:4], id[0:4]) != 1 {
		return errors.New("invalid id checksum")
	}
	return nil
}

func ObjectIdToPublicId(id bson.ObjectId) string {
	if id == "" {
		return ""
	}
	h := siphash.New(idHashKey)
	h.Write([]byte(id))
	return base64.URLEncoding.EncodeToString(append(h.Sum(nil)[0:4], []byte(id)...))
}

func PublicIdToObjectId(id string) (bson.ObjectId, error) {
	if id == "" {
		return "", nil
	}
	cypher, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}
	if err := validateId(cypher); err != nil {
		return "", err
	}
	return bson.ObjectId(cypher[4:16]), nil
}
