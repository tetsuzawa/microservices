package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
)

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{
		Error: errMessage,
		Code:  code,
	})
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(jsonError); err != nil {
		log.Println(err)
	}
}

const InternalServerErrorMessage = "Internal server error"

func InternalAPIError(w http.ResponseWriter) {
	APIError(w, InternalServerErrorMessage, http.StatusInternalServerError)
}


func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
