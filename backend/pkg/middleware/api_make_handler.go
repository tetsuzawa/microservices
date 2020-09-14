package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func ApiMakeHandler(next http.Handler, baseEndPoint string) http.Handler {
	re := fmt.Sprintf("^%s$", baseEndPoint)
	var apiValidPath = regexp.MustCompile(re)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		m := apiValidPath.FindStringSubmatch(r.URL.Path)
		if len(m) == 0 {
			w.WriteHeader(http.StatusNotFound)
			jsonError, err := json.Marshal(struct {
				Code  int    `json:"code"`
				Error string `json:"error"`
			}{
				http.StatusNotFound,
				"Not found",
			})
			if err != nil {
				log.Fatal(err)
			}
			if _, err := w.Write(jsonError); err != nil {
				log.Println(err)
			}
		}
		next.ServeHTTP(w, r)
	})
}
