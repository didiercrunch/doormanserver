package logpublisher

import (
	"log"

	"github.com/didiercrunch/doormanserver/shared"
	"gopkg.in/mgo.v2/bson"
)

type LogPublisher struct {
	Printer func(v ...interface{})
}

func (pub *LogPublisher) Emit(doormanId bson.ObjectId, doorman shared.AsDoormanUpdatePayloader) error {
	if data, err := doorman.AsDoormanUpdatePayload(); err != nil {
		return err
	} else {
		pub.Printer(doormanId.Hex(), " : ", string(data))
		return nil
	}
}

func (pub *LogPublisher) Init() error {
	pub.Printer = log.Println
	return nil
}
