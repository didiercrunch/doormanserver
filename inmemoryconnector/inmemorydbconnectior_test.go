package inmemoryconnector

import (
	"testing"

	"github.com/didiercrunch/doormanserver/doormen"
	"gopkg.in/mgo.v2/bson"
)

func createMockDoorman(name string) *doormen.DoormanDefinition {
	return doormen.QuickNewDoormanDefinition(name, 0.5, 0.5)
}

func TestNewInMemoryDBConnector(t *testing.T) {
	c := NewInMemoryDBConnector()
	if len(c.memory) != 0 {
		t.Error("length at this point should be 0")
	}
}

func TestSave(t *testing.T) {
	c := NewInMemoryDBConnector()
	if err := c.Save(nil); err != nil {
		t.Error(err)
	}
	if len(c.memory) != 1 {
		t.Fail()
	}
}

func TestDeleteDoorman(t *testing.T) {
	c := NewInMemoryDBConnector()
	c.Save(createMockDoorman("a"))
	c.Save(createMockDoorman("b"))
	c.Save(createMockDoorman("c"))

	if err := c.DeleteDoorman(c.memory[1].Id); err != nil {
		t.Error()
	}
	if c.memory[0].Name != "a" || c.memory[1].Name != "c" {
		t.Error()
	}

	if err := c.DeleteDoorman(c.memory[1].Id); err != nil {
		t.Error()
	}
	if c.memory[0].Name != "a" {
		t.Error()
	}

	if err := c.DeleteDoorman(c.memory[0].Id); err != nil {
		t.Error()
	}
	if len(c.memory) != 0 {
		t.Error()
	}
}

func TestGetAllDoormen(t *testing.T) {
	c := NewInMemoryDBConnector()
	c.Save(createMockDoorman("a"))
	c.Save(createMockDoorman("b"))
	c.Save(createMockDoorman("c"))

	i := 0
	expt := []string{"a", "b", "c"}
	for w := range c.GetAllDoormen() {
		if w.Name != expt[i] {
			t.Error(i)
		}
		i++
	}

}

func TestGetDoorman(t *testing.T) {
	c := NewInMemoryDBConnector()
	wab := createMockDoorman("a")
	c.Save(wab)

	if w, err := c.GetDoorman(wab.Id); err != nil {
		t.Error(err)
	} else if w.Name != "a" {
		t.Error()
	}

	if w, err := c.GetDoorman(createMockDoorman("b").Id); err != nil {
		t.Error("missing doorman is not an error")
	} else if w != nil {
		t.Error("should be nil")
	}
}

func TestExistsId(t *testing.T) {
	c := NewInMemoryDBConnector()
	w := createMockDoorman("a")
	c.Save(w)

	if !c.ExistsId(w.Id) {
		t.Error()
	}
	if c.ExistsId(bson.NewObjectId()) {
		t.Error()
	}

}
