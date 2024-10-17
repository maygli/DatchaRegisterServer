package commandserver

import (
	"datcha/authentification"
	"datcha/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type CommandServer struct {
	repository repository.IRepository
	authServer *authentification.AuthServer
}

func NewCommandServer(authServer *authentification.AuthServer, rep repository.IRepository) *CommandServer {
	server := CommandServer{}
	server.repository = rep
	server.authServer = authServer
	return &server
}

func (server *CommandServer) RegisterHandlers(r *mux.Router) {
	r.Handle("/command", server.authServer.JwtAuth(http.HandlerFunc(server.putCommandHandle))).Methods("PUT")
}
