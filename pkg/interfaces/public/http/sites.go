package http

import (
	"net/http"

	app "github.com/waffleboot/playstation_buy/pkg/application"
)

func sites(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	err := app.Query(search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(search))
}
