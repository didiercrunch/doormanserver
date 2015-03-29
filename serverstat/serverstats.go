package serverstat

import (
	"runtime"
	"strings"

	"github.com/didiercrunch/doormanserver/shared"
)

type m map[string]interface{}

var GO_VERSION = runtime.Version()

type ServerStat struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Go       string `json:"go_version"`

	MessageQueue string `json:"message_queue"`
	NSQ          m      `json:"nsq,omitempty"`
	Nanomsg      m      `json:"nano_msg,omitempty"`
}

func Get(p *shared.Params) *ServerStat {
	r := new(ServerStat)

	r.Hostname = p.Host
	r.Port = p.Port
	r.Go = GO_VERSION
	r.MessageQueue = strings.ToLower(p.MessageQueue)
	if p.UseNSQ() {
		r.NSQ = m{"ns_lookup_url": p.NSQ.NSQLookupdUrl}

	}
	if p.UseNanomsg() {
		r.Nanomsg = m{"url": p.NanoMsg.NanoMsgUrl}
	}
	return r
}
