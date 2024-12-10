package server

import (
	"encoding/json"
	"net/http"
)

func ToJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.Write([]byte("Error máximo"))
	}
}

func ToError(w http.ResponseWriter, data string, status int) {
	http.Error(w, data, status)
}
