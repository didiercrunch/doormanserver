package simpleapi

import (
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"path"
)

func JSONHendler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Type", mime.TypeByExtension(".json"))
		handler(w, request)
	}
}

type Endpoint struct {
	Pattern     string
	Description string
	Method      string
	HandlerFunc http.HandlerFunc
}

type SimpleApi struct {
	apiEndpoints []*Endpoint
	router       *mux.Router
	prefix       string
}

func New(prefix string, apiEndpoints ...*Endpoint) *SimpleApi {
	api := &SimpleApi{apiEndpoints, nil, prefix}
	api.initRouter()
	return api
}

func (api *SimpleApi) initRouter() {
	api.router = mux.NewRouter()
	for _, apiEndpoint := range api.apiEndpoints {
		api.router.Path(path.Join(api.prefix, apiEndpoint.Pattern)).
			Methods(apiEndpoint.Method).
			HandlerFunc(JSONHendler(apiEndpoint.HandlerFunc))
	}
}

func (api *SimpleApi) Set404Handler(handler http.HandlerFunc) {
	api.router.NotFoundHandler = http.HandlerFunc(handler)

}

func (api *SimpleApi) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	api.router.ServeHTTP(w, request)
}
