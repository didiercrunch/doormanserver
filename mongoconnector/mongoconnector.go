package mongoconnector

import (
	"github.com/didiercrunch/doormanserver/shared"
	"github.com/didiercrunch/doormanserver/doormen"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBConnector struct {
	Collection *mgo.Collection
}

func NewMongoDBConnector() (*MongoDBConnector, error) {
	if err := initMongo(); err != nil {
		return nil, err
	}

	c := new(MongoDBConnector)
	params := shared.GetParams()
	c.Collection = DB.C(params.Mongo.DoormanDefinitionCollection)
	err := c.Collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
	return c, err
}

func (c *MongoDBConnector) Exists(wab *doormen.DoormanDefinition) bool {
	err := c.Collection.Find(bson.M{"name": wab.Name}).One(&struct{}{})
	return err != mgo.ErrNotFound
}

func (c *MongoDBConnector) ExistsId(id bson.ObjectId) bool {
	err := c.Collection.FindId(id).One(&struct{}{})
	return err != mgo.ErrNotFound
}

func (c *MongoDBConnector) Save(wab *doormen.DoormanDefinition) error {
	return c.Collection.Insert(wab)
}

func (c *MongoDBConnector) GetAllDoormen() <-chan *doormen.DoormanId {
	ch := make(chan *doormen.DoormanId)
	go func() {
		iter := c.Collection.Find(nil).Iter()

		for res := new(doormen.DoormanId); iter.Next(res); res = new(doormen.DoormanId) {
			ch <- res
		}
		if err := iter.Err(); err != nil {
			ch <- &doormen.DoormanId{Id: "", Name: err.Error()}
		}
		close(ch)
	}()
	return ch
}

func (c *MongoDBConnector) GetDoorman(id bson.ObjectId) (*doormen.DoormanDefinition, error) {
	doorman := new(doormen.DoormanDefinition)
	if err := c.Collection.FindId(id).One(doorman); err == mgo.ErrNotFound {
		return nil, nil
	} else {
		return doorman, err
	}
}

func (c *MongoDBConnector) DeleteDoorman(id bson.ObjectId) error {
	return c.Collection.RemoveId(id)
}

func (c *MongoDBConnector) Update(wab *doormen.DoormanDefinition) error {
	return c.Collection.Update(&bson.M{"_id": wab.Id}, &bson.M{"$set": wab})
}
