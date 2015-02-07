package nanomsgpublisher

import (
	"errors"
	"log"

	"github.com/didiercrunch/doormanserver/shared"
	"github.com/gdamore/mangos"
	"github.com/gdamore/mangos/protocol/pub"
	"github.com/gdamore/mangos/transport/ipc"
	"github.com/gdamore/mangos/transport/tcp"
	"gopkg.in/mgo.v2/bson"
)

type NanoMsgPublisher struct {
	Url  string
	sock mangos.Socket
}

func (p *NanoMsgPublisher) Init() error {
	var err error
	if p.sock, err = pub.NewSocket(); err != nil {
		return errors.New("can't listen on pub socket:" + err.Error())
	}
	p.sock.AddTransport(ipc.NewTransport())
	p.sock.AddTransport(tcp.NewTransport())
	if err = p.sock.Listen(shared.GetParams().NanoMsg.NanoMsgUrl); err != nil {
		return errors.New("can't listen on pub socket: %s" + err.Error())
	}
	return nil

}

func (p *NanoMsgPublisher) Emit(doormanId bson.ObjectId, doorman shared.AsDoormanUpdatePayloader) error {
	log.Println("Emit for ", doormanId.Hex())
	if d, err := doorman.AsDoormanUpdatePayload(); err != nil {
		return err
	} else {
		return p.sock.Send(d)
	}
}
