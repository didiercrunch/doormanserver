package api

import (
	"encoding/json"
	"net/http"
)

func GetDoorman(w http.ResponseWriter, r *http.Request) {
	id := getDoormanIdFromRequest(w, r)
	if id == "" {
		Write404Error(w)
		return
	}
	doorman, err := conn.GetDoorman(id)
	if err != nil {
		Write500Error(w, err)
	} else if doorman == nil {
		Write404Error(w)
	} else {
		encoder := json.NewEncoder(w)
		encoder.Encode(doorman)
	}
}
