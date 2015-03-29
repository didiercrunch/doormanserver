package api

import (
	"errors"
	"log"

	"github.com/didiercrunch/doormanserver/doormen"
	"github.com/didiercrunch/doormanserver/inmemoryconnector"
	"github.com/didiercrunch/doormanserver/mongoconnector"
	"github.com/didiercrunch/doormanserver/shared"
	"gopkg.in/mgo.v2/bson"
)

type DBConnector interface {
	Save(wab *doormen.DoormanDefinition) error
	Update(wab *doormen.DoormanDefinition) error
	Exists(wab *doormen.DoormanDefinition) bool
	ExistsId(id bson.ObjectId) bool
	GetAllDoormen() <-chan *doormen.DoormanId
	GetDoorman(id bson.ObjectId) (*doormen.DoormanDefinition, error)
	DeleteDoorman(id bson.ObjectId) error
}

var conn DBConnector

func initMongoDbConnector() {
	var err error
	if conn, err = mongoconnector.NewMongoDBConnector(); err != nil {
		panic(err)
	}
}

func init() {
	if params := shared.GetParams(); params.UseMongoDb() {
		initMongoDbConnector()
		log.Println("using mongodb as database")
	} else {
		conn = inmemoryconnector.NewInMemoryDBConnector()
		log.Println("using inmemory as database")
	}
}

func CreateDoormanInDatabase(wab *doormen.DoormanDefinition) error {
	if conn.Exists(wab) {
		return errors.New("doormen already exists")
	}
	return conn.Save(wab)
}

func UpdateDoormanInDatabase(wab *doormen.DoormanDefinition) error {
	return conn.Update(wab)
}
