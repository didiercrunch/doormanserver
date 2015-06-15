package simpleapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"path"
)

type MiddleWare func(handler http.HandlerFunc) http.HandlerFunc

func JSONHandler(handler http.HandlerFunc) http.HandlerFunc {
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

func (ep *Endpoint) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"pattern":     path.Join(ep.Pattern),
		"description": ep.Description,
		"method":      ep.Method,
	}
	return json.Marshal(m)
}

type SimpleApi struct {
	apiEndpoints       []*Endpoint
	router             *mux.Router
	prefix             string
	middlewares        []MiddleWare
	documentationRoute string
}

func New(prefix string, apiEndpoints ...*Endpoint) *SimpleApi {
	api := &SimpleApi{apiEndpoints, nil, prefix, make([]MiddleWare, 0), ""}
	api.InitRouter()
	return api
}

func (api *SimpleApi) DocumentationHandler(w http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(w)
	payload := map[string]interface{}{"endpoints": api.apiEndpoints}
	if e := enc.Encode(payload); e != nil {
		fmt.Fprintln(w, e)
	}
}

func (api *SimpleApi) EnableDocumentation(route string) {
	if route == "" {
		return
	}
	api.documentationRoute = route
	api.router.Path(path.Join(api.prefix, route)).
		HandlerFunc(JSONHandler(api.DocumentationHandler))
}

func (api *SimpleApi) AddMiddlewares(middlewares ...MiddleWare) {
	api.middlewares = append(api.middlewares, middlewares...)
}

func (api *SimpleApi) applyMiddlewares(handler http.HandlerFunc) http.HandlerFunc {
	ret := handler
	for _, middleware := range api.middlewares {
		ret = middleware(ret)
	}
	return ret
}

func (api *SimpleApi) InitRouter() {
	api.router = mux.NewRouter()
	for _, apiEndpoint := range api.apiEndpoints {
		handler := JSONHandler(api.applyMiddlewares(apiEndpoint.HandlerFunc))
		api.router.Path(path.Join(api.prefix, apiEndpoint.Pattern)).
			Methods(apiEndpoint.Method).
			HandlerFunc(handler)
	}
	api.EnableDocumentation(api.documentationRoute)
}

func (api *SimpleApi) Set404Handler(handler http.HandlerFunc) {
	api.router.NotFoundHandler = http.HandlerFunc(handler)

}

func (api *SimpleApi) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	api.router.ServeHTTP(w, request)
}
