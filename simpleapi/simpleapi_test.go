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

func TestDocumentationHandler(t *testing.T) {
	api := New("/api", GetMockApiEndpoint("/foo/bar", "GET"))
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	api.DocumentationHandler(w, req)
	if code := w.Code; code != 200 {
		t.Error("bad status", code)
	}
	body := `{"endpoints":[{"description":"","method":"GET","pattern":"/foo/bar"}]}` + "\n"
	if body != w.Body.String() {
		t.Error()
	}
}

func TestEnableDocumentation(t *testing.T) {
	api := New("/api", GetMockApiEndpoint("/foo/bar", "GET"))
	api.EnableDocumentation("/doc")
	ts := httptest.NewServer(api)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/api/doc")
	if err != nil {
		t.Error(err)
		return
	}
	if res.StatusCode != 200 {
		t.Error("should have status 200 but got", res.StatusCode)
	}
}

func TestEnableDocumentationAndIniting(t *testing.T) {
	api := New("/api", GetMockApiEndpoint("/foo/bar", "GET"))
	api.EnableDocumentation("/doc")
	api.InitRouter()
	ts := httptest.NewServer(api)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/api/doc")
	if err != nil {
		t.Error(err)
		return
	}
	if res.StatusCode != 200 {
		t.Error("should have status 200 but got", res.StatusCode)
	}
}

func TestMiddleWares(t *testing.T) {
	api := New("/api", GetMockApiEndpoint("/foo", "GET"))
	middleWare1 := func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path += "/1"
			f(w, r)
		}
	}
	middleWare2 := func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path += "/2"
			f(w, r)
		}
	}
	api.AddMiddlewares(middleWare1, middleWare2)
	api.InitRouter()
	ts := httptest.NewServer(api)
	defer ts.Close()
	if rd, err := GET(ts, "/api/foo"); err != nil {
		t.Error(err)
	} else if rd.Path != "/api/foo/2/1" {
		t.Error(rd.Path)
	}

}
