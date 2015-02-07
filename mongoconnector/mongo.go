package mongoconnector

//var mongoClient int

import (
	"github.com/didiercrunch/doormanserver/shared"
	mgo "gopkg.in/mgo.v2"
)

var mongoClient *mgo.Session
var DB *mgo.Database

func initMongo() error {
	if mongoClient != nil {
		return nil
	}

	var err error
	mongoClient, err = mgo.Dial(shared.GetParams().Mongo.Url)
	if err != nil {
		return err
	}
	DB = mongoClient.DB(shared.GetParams().Mongo.DatabaseName)
	return nil
}
