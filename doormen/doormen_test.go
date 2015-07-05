package doormen

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/didiercrunch/doorman/shared"
	"gopkg.in/mgo.v2/bson"
)

var zero = big.NewRat(0, 1)

const epsilon float64 = 0.0001

func TestNewDoormanDefinition(t *testing.T) {
	dmd := NewDoormanDefinition("some_name")
	if dmd.Id == "" {
		t.Error("need id at all time")
	}
	if dmd.Name != "some_name" {
		t.Error("bad name")
	}
}

func TestValidate(t *testing.T) {
	dmd := &DoormanDefinition{Name: ""}
	if err := dmd.Validate("bob"); err.Error() != "name cannot be empty" {
		t.Error(err)
	}
	dmd = &DoormanDefinition{Name: "foo"}
	if err := dmd.Validate("bob"); err.Error() != "need at least 2 values" {
		t.Error(err)
	}
	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{nil}}
	if err := dmd.Validate("bob"); err.Error() != "need at least 2 values" {
		t.Error(err)
	}

	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", zero}, &DoormanValue{"foo", zero}}}
	if err := dmd.Validate("bob"); err.Error() != "doorman value names must be unique within a doorman" {
		t.Error(err)
	}

	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", zero}, &DoormanValue{"", zero}}}
	if err := dmd.Validate("bob"); err.Error() != "doorman value cannot be empty" {
		t.Error(err)
	}

	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", big.NewRat(-1, 4)}, &DoormanValue{"bar", big.NewRat(1, 2)}}}
	if err := dmd.Validate("bob"); err.Error() != "doorman value -1/4 is out of range" {
		t.Error(err)
	}

	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", big.NewRat(5, 2)}, &DoormanValue{"bar", big.NewRat(1, 2)}}}
	if err := dmd.Validate("bob"); err.Error() != "doorman value 5/2 is out of range" {
		t.Error(err)
	}

	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", big.NewRat(3, 4)}, &DoormanValue{"bar", big.NewRat(1, 2)}}}
	if err := dmd.Validate("bob"); err.Error() != "the sum of the probability must be 1" {
		t.Error(err)
	}

	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", big.NewRat(1, 5)}, &DoormanValue{"bar", big.NewRat(1, 5)}}}
	if err := dmd.Validate("bob"); err.Error() != "the sum of the probability must be 1" {
		t.Error(err)
	}
	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", big.NewRat(4, 5)}, &DoormanValue{"bar", big.NewRat(1, 5)}}}
	if err := dmd.Validate("bob"); err.Error() != "A doorman needs at least one owner specify as 'email'." {
		t.Error(err)
	}
	dmd = &DoormanDefinition{Name: "foo", Values: []*DoormanValue{&DoormanValue{"foo", big.NewRat(4, 5)}, &DoormanValue{"bar", big.NewRat(1, 5)}}, OwnerEmails: []string{"alice"}}
	if err := dmd.Validate("bob"); err.Error() != "'bob' is not allow to edit this doorman." {
		t.Error(err)
	}
}

func TestAsJson(t *testing.T) {
	id := &DoormanId{}
	if `{}` != id.AsJson() {
		t.Fail()
	}
	id.Name = "bob"
	if `{}` != id.AsJson() {
		t.Fail()
	}

	id.Id = bson.ObjectIdHex("54af464af6fb5e20c4000003")
	if `{"54af464af6fb5e20c4000003":"bob"}` != id.AsJson() {
		t.Error(id.AsJson())
	}
}

func TestAsDoormanUpdatePayload(t *testing.T) {
	wd := &DoormanDefinition{}
	wd.Id = bson.NewObjectId()
	wd.Values = []*DoormanValue{&DoormanValue{"C", big.NewRat(1, 4)}, &DoormanValue{"T1", big.NewRat(3, 4)}}
	data, err := wd.AsDoormanUpdatePayload()
	if err != nil {
		t.Error(err)
		return
	}
	wu := &shared.DoormanUpdater{}
	if err := json.Unmarshal(data, wu); err != nil {
		t.Error(err)
	}
	if wd.Id.Hex() != wu.Id {
		t.Error("bad id", wu.Id)
	}
	if len(wu.Probabilities) != 2 {
		t.Error("bad probability length")
	} else if wu.Probabilities[0].Cmp(big.NewRat(1, 4)) != 0 {
		t.Error("bad first probability")
	} else if wu.Probabilities[1].Cmp(big.NewRat(3, 4)) != 0 {
		t.Error("bad second probability")
	} else if wu.Timestamp == 0 {
		t.Error("no time specify")
	} else if wu.Timestamp < 1421537442 {
		t.Error("timestamp looks weird")
	}
}

func TestQuickNewDoormanDefinition(t *testing.T) {
	dmd := QuickNewDoormanDefinition("foo", "0.75", "0.25")
	if dmd.Name != "foo" {
		t.Error()
	}

	if len(dmd.Values) != 2 {
		t.Error()
	}

	if dmd.Values[0].Name != "T0" || dmd.Values[0].Probability.Cmp(big.NewRat(3, 4)) != 0 {
		t.Error()
	}

	if dmd.Values[1].Name != "T1" || dmd.Values[1].Probability.Cmp(big.NewRat(1, 4)) != 0 {
		t.Error()
	}
	if len(dmd.OwnerEmails) != 1 && dmd.OwnerEmails[0] != "natasha@bigtits.com" {
		t.Error()
	}
}

func TestAsWriteAccess(t *testing.T) {
	dmd := &DoormanDefinition{OwnerEmails: []string{"natasha@bigtits.com"}}
	testCases := map[string]bool{
		"natasha@bigtits.com": true,
		"natasha@bigtis.net":  false,
	}
	for email, hasAccess := range testCases {
		if hasAccess != dmd.AsWriteAccess(email) {
			t.Error(email)
		}
	}

}
