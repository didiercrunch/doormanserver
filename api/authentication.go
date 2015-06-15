package api

import (
	"github.com/gorilla/context"
	"net/http"
)

const userKey = "user"

func GetUser(req *http.Request) string {
	if rv := context.Get(req, userKey); rv != nil {
		return rv.(string)
	}
	return ""
}

func SetUser(req *http.Request, user string) {
	context.Set(req, userKey, user)
}

func AuthenticationMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		user := query.Get("user")
		SetUser(req, user)
		h(w, req)
		context.Clear(req)
	}
}
