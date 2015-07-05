package shared

import (
	"io/ioutil"
	"log"
	"strings"

	goyaml "gopkg.in/yaml.v2"
)

var PARAM_FILE string = "params.yml"

var params *Params

type MongoParams struct {
	Url                        string `yaml:"url,omitempty"`
	DatabaseName               string `yaml:"database_name,omitempty"`
	DoormanDefinitionCollection string `yaml:"doorman_definition_collection,omitempty"`
}

type NSQParams struct {
	NSQLookupdUrl string `yaml:"nsq_lookupd_url,omitempty"`
}

type NanoMsgParams struct {
	NanoMsgUrl string `yaml:"nanomsg_url,omitempty"`
}

type Params struct {
	Port    int            `yaml:"port,omitempty"`
	Host    string         `yaml:"host,omitempty"`
	Mongo   *MongoParams   `yaml:"mongo,omitempty"`
	NSQ     *NSQParams     `yaml:"nsq,omitempty"`
	NanoMsg *NanoMsgParams `yaml:"nanomsg,omitempty"`

	Database     string `yaml:"database,omitempty"`
	MessageQueue string `yaml:"message_queue,omitempty"`
}

func (p *Params) UseMongoDb() bool {
	db := strings.ToLower(p.Database)
	return db == "mongo" || db == "mongodb"
}

func (p *Params) UseNSQ() bool {
	m := strings.ToLower(p.MessageQueue)
	return m == "nsq"
}

func (p *Params) UseNanomsg() bool {
	m := strings.ToLower(p.MessageQueue)
	return m == "nanomsg"
}

func (p *Params) InsertDefaultValues() *Params {
	if p.Host == "" {
		p.Host = "localhost"
	}
	if p.Port == 0 {
		p.Port = 1999
	}
	if p.Mongo == nil {
		p.Mongo = new(MongoParams)
	}
	if p.NSQ == nil {
		p.NSQ = new(NSQParams)
	}
	return p
}

func LoadFromYamlFile(fileName string) (*Params, error) {
	if !exists(fileName) {
		return new(Params).InsertDefaultValues(), nil
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return LoadFromYamlData(data)
}

func LoadFromYamlData(data []byte) (*Params, error) {
	var err error
	p := new(Params)
	if err = goyaml.Unmarshal(data, p); err != nil {
		return nil, err
	}
	return p.InsertDefaultValues(), nil
}

func GetParams() *Params {
	return params
}

func init() {
	var err error
	if params, err = LoadFromYamlFile(PARAM_FILE); err != nil {
		log.Panic(err)
	}
}
