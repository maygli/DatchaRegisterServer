package channelserver

import (
	"datcha/authentification"
	"datcha/repository"
	"datcha/servercommon"
	"net/http"

	"github.com/gorilla/mux"
)

type ChannelServer struct {
	repository repository.IRepository
	authServer *authentification.AuthServer
}

func NewChannelServer(authServer *authentification.AuthServer, rep repository.IRepository) *ChannelServer {
	server := ChannelServer{}
	server.repository = rep
	server.authServer = authServer
	return &server
}

func (server *ChannelServer) RegisterHandlers(r *mux.Router) {
	r.Handle("/channel/data/{"+servercommon.CHANNEL_ID_KEY+"}", server.authServer.JwtAuth(http.HandlerFunc(server.getChannelDataHandle))).Methods("GET")
}
