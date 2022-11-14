package http_utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ApplicationJson() (string, string) {
	return "Content-Type", "application/json"
}

func WriteJsonBody(respWriter http.ResponseWriter, statusCode int, obj interface{}) {
	// convert object to JSON
	data, err := json.Marshal(obj)

	if err != nil {
		log.Printf("Can't properly marshall object to JSON: %v, error: %s\n", obj, err.Error())
		respWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write status and response body
	respWriter.WriteHeader(statusCode)
	_, writeErr := respWriter.Write(data)
	if writeErr != nil {
		log.Printf("Can't properly write response: %s", writeErr.Error())
		return
	}
}
