package api

import (
	"mime"
	"net/http"
)

func GetDormanState(w http.ResponseWriter, request *http.Request) {
	id := getDoormanIdFromRequest(w, request)
	if id == "" {
		Write404Error(w)
		return
	}
	doorman, err := conn.GetDoorman(id)
	if err != nil {
		Write500Error(w, err)
		return
	} else if doorman == nil {
		Write404Error(w)
		return
	}

	if payload, err := doorman.AsDoormanUpdatePayload(); err != nil {
		Write500Error(w, err)
	} else {
		w.Header().Add("Content-Type", mime.TypeByExtension(".json"))
		w.Write(payload)
	}
}
