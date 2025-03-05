package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSON takes in a dto and writes it to the response
func WriteJSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		log.Println("error marshaling json: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(js)
	if err != nil {
		log.Println("error writing json: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
