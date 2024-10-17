package pageserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandlers(theMux *mux.Router) {
	theMux.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
}
