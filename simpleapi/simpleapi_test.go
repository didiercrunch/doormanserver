package simpleapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type requestData struct {
	Path     string
	Method   string
	Response *http.Response
}

func createRequestData(r *http.Request) *requestData {
	return &requestData{r.URL.Path, r.Method, nil}

}

func GetHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		encoder.Encode(createRequestData(r))
	}
}

func GetMockApiEndpoint(pattern, method string) *Endpoint {
	return &Endpoint{Pattern: pattern, Method: method, HandlerFunc: GetHandlerFunc()}
}

func GET(ts *httptest.Server, path string) (*requestData, error) {
	res, err := http.Get(ts.URL + path)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	ret := new(requestData)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	ret.Response = res
	return ret, nil
}

func POST(ts *httptest.Server, path, data string) (*requestData, error) {
	res, err := http.Post(ts.URL+path, "", bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	ret := new(requestData)
	return ret, json.Unmarshal(body, ret)
}

func TestServeGet(t *testing.T) {
	api := New("/", GetMockApiEndpoint("/", "GET"))
	ts := httptest.NewServer(api)
	defer ts.Close()

	if rd, err := GET(ts, ""); err != nil {
		t.Error(err)
	} else if rd.Path != "/" {
		t.Error(rd)
	}
}

func TestServeGetSpecificPath(t *testing.T) {
	api := New("/", GetMockApiEndpoint("/some/path", "GET"))
	ts := httptest.NewServer(api)
	defer ts.Close()

	if rd, err := GET(ts, "/some/path"); err != nil {
		t.Error(err)
	} else if rd.Path != "/some/path" {
		t.Error(rd)
	}
}

func TestServePOSTSpecificPath(t *testing.T) {
	api := New("/", GetMockApiEndpoint("/some/path", "POST"))
	ts := httptest.NewServer(api)
	defer ts.Close()

	if rd, err := POST(ts, "/some/path", ""); err != nil {
		t.Error(err)
	} else if rd.Path != "/some/path" || rd.Method != "POST" {
		t.Error(rd)
	}
}

func Test404(t *testing.T) {
	api := New("/")
	api.Set404Handler(GetHandlerFunc())
	ts := httptest.NewServer(api)
	defer ts.Close()

	if rd, err := GET(ts, "/some/path"); err != nil {
		t.Error(err)
	} else if rd.Path != "/some/path" || rd.Method != "GET" {
		t.Error(rd)
	}
}

func TestHeaders(t *testing.T) {
	api := New("/", GetMockApiEndpoint("/", "GET"))
	ts := httptest.NewServer(api)
	defer ts.Close()

	if rd, err := GET(ts, ""); err != nil {
		t.Error(err)
	} else if rd.Response.Header.Get("Content-Type") != "application/json" {
		t.Error(rd)
	}
}

func TestServePrefix(t *testing.T) {
	api := New("/api", GetMockApiEndpoint("/foo/bar", "GET"))
	ts := httptest.NewServer(api)
	defer ts.Close()

	if rd, err := GET(ts, "/api/foo/bar"); err != nil {
		t.Error(err)
	} else if rd.Path != "/api/foo/bar" {
		t.Error(rd)
	}
}
