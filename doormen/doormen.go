package doormen

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math/big"
	"strconv"
	"time"

	client "github.com/didiercrunch/doorman/shared"
	"github.com/didiercrunch/doormanserver/shared"
)

type DoormanValue struct {
	Name        string   `json:"name"`
	Probability *big.Rat `json:"probability"`
}

type DoormanId struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}

func (wId *DoormanId) AsJson() string {
	if wId.Id == "" {
		return `{}`
	}
	return fmt.Sprintf(`{"%v":"%v"}`, wId.Id.Hex(), wId.Name)
}

type DoormanDefinition struct {
	Id          bson.ObjectId   `bson:"_id"`
	Name        string          `bson:"name"`
	Values      []*DoormanValue `bson:"values"`
	OwnerEmails []string        `bson:"emails"`
}

func NewDoormanDefinition(name string, ownerEmails ...string) *DoormanDefinition {
	return &DoormanDefinition{Id: bson.NewObjectId(), Name: name, OwnerEmails: ownerEmails}
}

// a way to quickly create doormen in tests.  do not use in real code.  will
// panic if a doorman is wrongly specify
func QuickNewDoormanDefinition(name string, probs ...string) *DoormanDefinition {
	ret := NewDoormanDefinition(name, "natasha@bigtits.com")
	for i, prob := range probs {
		r := new(big.Rat)
		_, err := fmt.Sscan(prob, r)
		if err != nil {
			panic(err)
		}
		ret.Values = append(ret.Values, &DoormanValue{"T" + strconv.Itoa(i), r})
	}
	if err := ret.Validate("natasha@bigtits.com"); err != nil {
		panic(err)
	}
	return ret
}

func (dmd *DoormanDefinition) ValidateDoormanValueNames() error {
	m := make(map[string]bool)
	for _, w := range dmd.Values {
		if name := w.Name; name == "" {
			return errors.New("doorman value cannot be empty")
		} else {
			m[name] = true
		}
	}
	if len(m) != len(dmd.Values) {
		return errors.New("doorman value names must be unique within a doorman")
	}
	return nil
}

func (dmd *DoormanDefinition) AsWriteAccess(email string) bool {
	for _, email_ := range dmd.OwnerEmails {
		if email_ == email {
			return true
		}
	}
	return false
}

func (dmd *DoormanDefinition) ValidateDoormanValueProbabilities() error {
	sum, zero, one := big.NewRat(0, 1), big.NewRat(0, 1), big.NewRat(1, 1)
	for _, w := range dmd.Values {

		if p := w.Probability; p.Cmp(zero) < 0 || p.Cmp(one) > 0 {
			return errors.New(fmt.Sprintf("doorman value %v is out of range", p))
		} else {
			sum = new(big.Rat).Add(sum, p)
		}
	}
	if sum.Cmp(one) != 0 {
		return errors.New("the sum of the probability must be 1")
	}

	return nil
}

func (dmd *DoormanDefinition) Validate(author string) (err error) {
	switch {
	case dmd.Name == "":
		err = errors.New("name cannot be empty")
	case len(dmd.Values) < 2:
		err = errors.New("need at least 2 values")
	case dmd.ValidateDoormanValueNames() != nil:
		err = dmd.ValidateDoormanValueNames()
	case dmd.ValidateDoormanValueProbabilities() != nil:
		err = dmd.ValidateDoormanValueProbabilities()
	case len(dmd.OwnerEmails) < 1:
		err = errors.New("A doorman needs at least one owner specify as 'email'.")
	case !dmd.AsWriteAccess(author):
		err = errors.New(fmt.Sprintf("'%s' is not allow to edit this doorman.", author))
	}
	return
}

func (dmd *DoormanDefinition) CanBeUpdatedBy(dmd2 *DoormanDefinition) error {
	if dmd.Name != dmd2.Name {
		return errors.New("cannot change doorman name")
	}
	if len(dmd.Values) != len(dmd2.Values) {
		return errors.New("cannot change the number of cases in doorman")
	}
	return nil
}

func (dmd *DoormanDefinition) AsDoormanUpdatePayload() ([]byte, error) {
	wu := &client.DoormanUpdater{Id: shared.ObjectIdToPublicId(dmd.Id)}
	wu.Timestamp = time.Now().Unix()
	wu.Probabilities = make([]*big.Rat, len(dmd.Values))
	for i, value := range dmd.Values {
		wu.Probabilities[i] = value.Probability
	}
	return json.Marshal(wu)
}
