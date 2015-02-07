package inmemoryconnector

import (
	"errors"

	"github.com/didiercrunch/doormanserver/doormen"
	"gopkg.in/mgo.v2/bson"
)

type InMemoryDBConnector struct {
	memory []*doormen.DoormanDefinition
}

func NewInMemoryDBConnector() *InMemoryDBConnector {
	return &InMemoryDBConnector{make([]*doormen.DoormanDefinition, 0)}
}

func (c *InMemoryDBConnector) Save(dmd *doormen.DoormanDefinition) error {
	c.memory = append(c.memory, dmd)
	return nil
}

func (c *InMemoryDBConnector) Exists(dmd *doormen.DoormanDefinition) bool {
	for _, w := range c.memory {
		if w.Name == dmd.Name {
			return true
		}
	}
	return false
}

func (c *InMemoryDBConnector) ExistsId(id bson.ObjectId) bool {
	for _, w := range c.memory {
		if w.Id == id {
			return true
		}
	}
	return false

}

func (c *InMemoryDBConnector) GetAllDoormen() <-chan *doormen.DoormanId {
	ch := make(chan *doormen.DoormanId)
	go func() {
		for _, s := range c.memory {
			ch <- &doormen.DoormanId{Name: s.Name, Id: s.Id}
		}
		close(ch)
	}()
	return ch

}

func (c *InMemoryDBConnector) GetDoorman(id bson.ObjectId) (*doormen.DoormanDefinition, error) {
	for _, doorman := range c.memory {
		if doorman.Id == id {
			return doorman, nil
		}
	}
	return nil, nil
}

func (c *InMemoryDBConnector) DeleteDoorman(id bson.ObjectId) error {
	toRemove := -1
	for i, doorman := range c.memory {
		if doorman.Id == id {
			toRemove = i
			break
		}
	}
	c.memory = append(c.memory[0:toRemove], c.memory[toRemove+1:len(c.memory)]...)
	return nil
}

func (c *InMemoryDBConnector) Update(dmd *doormen.DoormanDefinition) error {
	for i, w := range c.memory {
		if w.Id == dmd.Id {
			c.memory[i] = dmd
			return nil
		}
	}
	return errors.New("nothing to update")
}
