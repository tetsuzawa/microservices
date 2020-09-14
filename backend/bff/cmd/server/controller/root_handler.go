package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if head == "" && r.Method == http.MethodGet {
		Index(w)
	}
}

func Index(w http.ResponseWriter) {
	js, err := json.Marshal(struct{ Result string `json:"result"` }{"root handler"})
	if err != nil {
		log.Println(err)
		InternalAPIError(w)
	}
	if _, err := w.Write(js); err != nil {
		log.Println(err)
	}
}
