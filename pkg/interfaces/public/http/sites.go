package http

import (
	"net/http"
)

func sites(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	w.Write([]byte(search))
}
