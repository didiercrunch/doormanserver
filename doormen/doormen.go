package doormen

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/didiercrunch/doorman/shared"
	"gopkg.in/mgo.v2/bson"
)

const epsilon float64 = 0.0001

type DoormanValue struct {
	Name        string  `json:"name" bson:"name"`
	Probability float64 `json:"probability" bson:"probability"`
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
	Id     bson.ObjectId   `json:"id,omitempty" bson:"_id"`
	Name   string          `json:"name" bson:"name"`
	Values []*DoormanValue `json:"values" bson:"values"`
}

func NewDoormanDefinition(name string) *DoormanDefinition {
	return &DoormanDefinition{Id: bson.NewObjectId(), Name: name}
}

// a way to quickly create doormen in tests.  do not use in real code.  will
// panic if a doorman is wrongly specify
func QuickNewDoormanDefinition(name string, probs ...float64) *DoormanDefinition {
	ret := NewDoormanDefinition(name)
	for i, prob := range probs {
		ret.Values = append(ret.Values, &DoormanValue{"T" + strconv.Itoa(i), prob})
	}
	if err := ret.Validate(); err != nil {
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

func (dmd *DoormanDefinition) ValidateDoormanValueProbabilities() error {
	var sum float64 = 0
	for _, w := range dmd.Values {

		if p := w.Probability; p < 0-epsilon || p > 1+epsilon {
			return errors.New(fmt.Sprintf("doorman value %v is out of range", p))
		} else {
			sum += p
		}
	}
	if sum < 1-epsilon || sum > 1+epsilon {
		return errors.New("the sum of the probability must be 1.0")
	}
	return nil
}

func (dmd *DoormanDefinition) Validate() (err error) {
	switch {
	case dmd.Name == "":
		err = errors.New("name cannot be empty")
	case len(dmd.Values) < 2:
		err = errors.New("need at least 2 values")
	case dmd.ValidateDoormanValueNames() != nil:
		err = dmd.ValidateDoormanValueNames()
	case dmd.ValidateDoormanValueProbabilities() != nil:
		err = dmd.ValidateDoormanValueProbabilities()
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
	wu := &shared.DoormanUpdater{Id: dmd.Id.Hex()}
	wu.Timestamp = time.Now().Unix()
	wu.Probabilities = make([]float64, len(dmd.Values))
	for i, value := range dmd.Values {
		wu.Probabilities[i] = value.Probability
	}
	return json.Marshal(wu)
}
