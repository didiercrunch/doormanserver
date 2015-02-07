package shared

import (
	"testing"
)

const yaml_params = `
---
port: 9999
host: 0.0.0.0  #localhost
# all the params relative to mongodb
mongo:
    # the url where to get mongodb
    url: localhost:1234

    # the name of the database to use
    database_name: doorman

    # the name of the collection where all the street can be found
    doorman_definition_collection: doorman_definitions

#  all the parameters relative to nsq message-queue
nsq:
    # the url where to find nsq lookupd url
    nsq_lookupd_url: 127.0.0.1:4161

#  all the parameters relative to nanomsg message-queue
nanomsg:
    nanomsg_url: ipc:///tmp/doorman.ipc


`

func TestLoadFromYamlData(t *testing.T) {
	p, err := LoadFromYamlData([]byte(yaml_params))
	if err != nil {
		t.Error(err)
		return
	}
	if p.Port != 9999 {
		t.Error("port")
	}
	if p.Host != "0.0.0.0" {
		t.Error("host")
	}
}

func TestMongoParams(t *testing.T) {
	p, err := LoadFromYamlData([]byte(yaml_params))
	if err != nil {
		t.Error(err)
		return
	}
	if p.Mongo.Url != "localhost:1234" {
		t.Error("mongo.url")
	}
	if p.Mongo.DatabaseName != "doorman" {
		t.Error("mongo.DatabaseName")
	}
	if p.Mongo.DoormanDefinitionCollection != "doorman_definitions" {
		t.Error("Mongo.DoormanDefinitionCollection", p.Mongo.DoormanDefinitionCollection)
	}
}

func TestNSQParams(t *testing.T) {
	p, err := LoadFromYamlData([]byte(yaml_params))
	if err != nil {
		t.Error(err)
		return
	}
	if p.NSQ.NSQLookupdUrl != "127.0.0.1:4161" {
		t.Error("Nsq.NSQLookupdUrl", p.NSQ.NSQLookupdUrl)
	}
}

func TestNanoMsgParams(t *testing.T) {
	p, err := LoadFromYamlData([]byte(yaml_params))
	if err != nil {
		t.Error(err)
		return
	}
	if p.NanoMsg.NanoMsgUrl != "ipc:///tmp/doorman.ipc" {
		t.Error("NanoMsgParams.NanoMsgUrl", p.NanoMsg.NanoMsgUrl)
	}
}

func TestMinimalParams(t *testing.T) {
	params := "---\n"
	p, err := LoadFromYamlData([]byte(params))
	if err != nil {
		t.Error(err)
	}
	if p.Mongo.Url != "" {
		t.Error()
	}
	if p.Port != 1999 {
		t.Error("bad default port", p.Port)
	}
	if p.Host != "localhost" {
		t.Error("bad default host")
	}
	if p.MessageQueue != "" {
		t.Error("message queue")
	}

	if p.Database != "" {
		t.Error("database")
	}
}

func TestUseMongoDb(t *testing.T) {
	p := &Params{}
	if p.Database = ""; p.UseMongoDb() {
		t.Error()
	}
	if p.Database = "mOnGo"; !p.UseMongoDb() {
		t.Error()
	}
	if p.Database = "monGoDb"; !p.UseMongoDb() {
		t.Error()
	}
}

func TestUseNSQ(t *testing.T) {
	p := &Params{}
	if p.MessageQueue = ""; p.UseNSQ() {
		t.Error()
	}
	if p.MessageQueue = "nSq"; !p.UseNSQ() {
		t.Error()
	}
}

func TestUseNanoMsg(t *testing.T) {
	p := &Params{}
	if p.MessageQueue = ""; p.UseNanomsg() {
		t.Error()
	}
	if p.MessageQueue = "nAnOmSg"; !p.UseNanomsg() {
		t.Error()
	}
}
