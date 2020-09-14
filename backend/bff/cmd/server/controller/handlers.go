package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Result string `json:"result"`
	}{"OK"}
	js, err := json.Marshal(res)
	if err != nil {
		APIError(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err := w.Write(js); err != nil {
		log.Println(err)
	}
}
