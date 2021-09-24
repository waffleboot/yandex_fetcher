package http

import (
	"net/http"
)

func (e *Endpoint) sites(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if err := e.service.Query(search); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(search))
}
