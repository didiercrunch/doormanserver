package doormen

import (
	"encoding/json"
	"github.com/didiercrunch/doormanserver/shared"
	"gopkg.in/mgo.v2/bson"
	"math/big"
)

type doormanValue struct {
	Name        string
	Probability string
}

func (val *DoormanValue) GetBSON() (interface{}, error) {
	ret := &doormanValue{val.Name, ""}
	txt, err := val.Probability.MarshalText()
	if err != nil {
		return nil, err
	}
	ret.Probability = string(txt)
	return ret, nil
}

func (val *DoormanValue) SetBSON(raw bson.Raw) error {
	v := new(doormanValue)
	if err := raw.Unmarshal(v); err != nil {
		return err
	}
	val.Name = v.Name
	val.Probability = new(big.Rat)
	return val.Probability.UnmarshalText([]byte(v.Probability))
}

type jsonDoormanDefinition struct {
	Id          string          `json:"id,omitempty"`
	Name        string          `json:"name"`
	Values      []*DoormanValue `json:"values"`
	OwnerEmails []string        `json:"emails"`
}

func (dm *DoormanDefinition) MarshalJSON() ([]byte, error) {
	jdm := new(jsonDoormanDefinition)
	jdm.Name = dm.Name
	jdm.OwnerEmails = dm.OwnerEmails
	jdm.Values = dm.Values
	jdm.Id = shared.ObjectIdToPublicId(dm.Id)
	return json.Marshal(jdm)
}

func (dm *DoormanDefinition) UnmarshalJSON(data []byte) error {
	jsm := new(jsonDoormanDefinition)
	var err error
	if err = json.Unmarshal(data, jsm); err != nil {
		return err
	}
	if dm.Id, err = shared.PublicIdToObjectId(jsm.Id); err != nil {
		return err
	}
	dm.Name = jsm.Name
	dm.Values = jsm.Values
	dm.OwnerEmails = jsm.OwnerEmails
	return nil
}

type jsonDoormanId struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (did *DoormanId) MarshalJSON() ([]byte, error) {
	jdid := new(jsonDoormanId)
	jdid.Name = did.Name
	jdid.Id = shared.ObjectIdToPublicId(did.Id)
	return json.Marshal(jdid)

}

func (did *DoormanId) UnmarshalJSON(data []byte) error {
	jdid := new(jsonDoormanId)
	if err := json.Unmarshal(data, jdid); err != nil {
		return err
	}
	var err error
	if did.Id, err = shared.PublicIdToObjectId(jdid.Id); err != nil {
		return err
	}
	did.Name = jdid.Name
	return nil
}
