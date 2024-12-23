// Code get from https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c
package notifyservice

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

func NewNotifyService() *Broker {
	// Instantiate a broker
	broker := &Broker{
		Notifier:       make(chan []byte),
		clients:        make(map[chan []byte]bool),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return broker
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	/*	user, ok := req.Context().Value(servercommon.USER_CONTEXT_KEY).(datamodel.User)
		if !ok {
			log.Println("Notifieer: Context doesn't contains user information")
			http.Error(rw, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
			return
		}
		log.Println("Notifier with user=" + user.Name)*/
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	// notify := rw.(http.CloseNotifier).CloseNotify()
	notify := req.Context().Done()

	for {
		isBreak := false
		select {
		case event := <-messageChan:
			// Write to the ResponseWriter
			// Server Sent Events compatible
			fmt.Fprintf(rw, "data: %s\n\n", event)

			// Flush the data immediatly instead of buffering it for later.
			flusher.Flush()

		case <-notify:
			isBreak = true
		}
		if isBreak {
			break
		}
	}
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			slog.Debug(fmt.Sprintf("Client added. %d registered clients", len(broker.clients)))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			slog.Debug(fmt.Sprintf("Removed client. %d registered clients", len(broker.clients)))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (broker Broker) Notify(event []byte) {
	broker.Notifier <- event
}
