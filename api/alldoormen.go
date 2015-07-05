package api

import (
	"fmt"
	"net/http"
)

func GetAllDoormen(w http.ResponseWriter, request *http.Request) {
	doormenIds := conn.GetAllDoormen()
	buff := <-doormenIds
	if buff == nil {
		fmt.Fprint(w, `{"doormen" :[]}`)
		return
	}
	fmt.Fprint(w, `{"doormen" :[`)
	for doormanId := range doormenIds {
		fmt.Fprint(w, getDoormenIdsAsJson(buff)+`,`)
		buff = doormanId
	}
	fmt.Fprint(w, getDoormenIdsAsJson(buff), "]}")
}
