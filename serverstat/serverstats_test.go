package serverstat

import (
	"encoding/json"
	"reflect"
	"testing"
)
import "github.com/didiercrunch/doormanserver/shared"

func deserializeServerStats(data []byte) *ServerStat {
	ret := new(ServerStat)
	if err := json.Unmarshal(data, ret); err != nil {
		panic("cannot deserialize data")
	}
	return ret
}

func ServerStatsEqual(a, b *ServerStat) bool {
	return reflect.DeepEqual(a, b)
}

func TestGetTrivial(t *testing.T) {
	p := new(shared.Params)
	expected := &ServerStat{Go: GO_VERSION}
	if s := Get(p); !ServerStatsEqual(s, expected) {
		t.Error(s)
	}
}
