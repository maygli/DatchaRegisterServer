package notifyservice

import (
	"net/http"
)

func (broker *Broker) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /notify", http.HandlerFunc(broker.ServeHTTP))
}
