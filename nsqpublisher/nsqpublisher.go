package nsqpublisher

import (
	"github.com/bitly/go-nsq"
	"github.com/didiercrunch/doormanserver/shared"
	"gopkg.in/mgo.v2/bson"
)

type NsqPublisher struct {
	Config   *nsq.Config
	Producer *nsq.Producer
}

func (pub *NsqPublisher) InitConfig() {
	pub.Config = nsq.NewConfig()
}

func (pub *NsqPublisher) Init() error {
	var err error
	pub.InitConfig()
	pub.Producer, err = nsq.NewProducer(shared.GetParams().NSQ.NSQLookupdUrl, pub.Config)
	return err

}

func (pub *NsqPublisher) Stop() error {
	pub.Producer.Stop()
	return nil
}

func (pub *NsqPublisher) Emit(doormanId bson.ObjectId, doorman shared.AsDoormanUpdatePayloader) error {
	topic := doormanId.Hex()
	if data, err := doorman.AsDoormanUpdatePayload(); err != nil {
		return err
	} else {
		return pub.Producer.Publish(topic, data)
	}
}
