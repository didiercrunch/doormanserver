package main

import (
	"log"

	"github.com/didiercrunch/doormanserver/logpublisher"
	"github.com/didiercrunch/doormanserver/nanomsgpublisher"
	"github.com/didiercrunch/doormanserver/nsqpublisher"
	"github.com/didiercrunch/doormanserver/shared"
	"gopkg.in/mgo.v2/bson"
)

type Publisher interface {
	Emit(doormanId bson.ObjectId, doorman shared.AsDoormanUpdatePayloader) error
	Init() error
}

var publisher Publisher

func initNSQPublisher() {
	publisher = new(nsqpublisher.NsqPublisher)
	if err := publisher.Init(); err != nil {
		log.Panic(err)
		return
	}
}

func initNanoMsgPublisher() {
	publisher = new(nanomsgpublisher.NanoMsgPublisher)
	if err := publisher.Init(); err != nil {
		log.Panic(err)
		return
	}
}

func initPublisher(params *shared.Params) {
	if params.UseNSQ() {
		initNSQPublisher()
		log.Println("using nsq as message queue")

	} else if params.UseNanomsg() {
		initNanoMsgPublisher()
		log.Println("using nanomsg as message queue")
	} else {
		publisher = &logpublisher.LogPublisher{log.Println}
		log.Println("using stdout as message queue")
	}

}
func init() {
	initPublisher(shared.GetParams())
}
