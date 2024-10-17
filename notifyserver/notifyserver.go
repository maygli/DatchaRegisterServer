package notifyserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (broker *Broker) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/notify", http.HandlerFunc(broker.ServeHTTP)).Methods("GET")
}
