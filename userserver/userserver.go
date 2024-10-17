package userserver

import (
	"datcha/authentification"
	"datcha/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type UserServer struct {
	repository repository.IRepository
	authServer authentification.AuthServer
}

func (server *UserServer) RegisterHandlers(r *mux.Router) {
	r.Handle("/user/profile", server.authServer.JwtAuth(http.HandlerFunc(server.getUserProfileHandle))).Methods(http.MethodGet)
	r.Handle("/user/profile", server.authServer.JwtAuth(http.HandlerFunc(server.updateUserProfileHandle))).Methods(http.MethodPut)
}
