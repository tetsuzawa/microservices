package main

import (
	"net/http"

	"github.com/tetsuzawa/microservices/backend/bff/cmd/server/controller"
)

type ServerHandler struct {
	RootHandler *controller.RootHandler
}

func (h *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = controller.ShiftPath(r.URL.Path)

	switch head {
	case "":
		h.RootHandler.ServeHTTP(w, r)
		return
		//case "users":
		//	h.UserHandler.ServeHTTP(w, r)
		//	return
	}
}
