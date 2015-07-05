package doormen

import (
	"gopkg.in/mgo.v2/bson"
	"math/big"
)

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
