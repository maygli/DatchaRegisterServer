package authentification

import (
	"datcha/servercommon"
	"log"
	"net/http"
	"strings"
)

func (server *AuthServer) refreshPostHandle(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(AUTHORIZATION_HEADER)
	headerData := strings.Split(authHeader, BEARER)
	if len(headerData) < 2 {
		log.Println("Authorisation header is not found or contains invalid information")
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	tokenData := strings.TrimSpace(headerData[1])
	userId, err := server.getUserIdFromToken(tokenData, REFRESH_SUBJECT, server.RefreshSecretKey)
	if err != nil {
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	_, err = server.repository.FindUser(userId)
	if err != nil {
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	server.sendAuthCoockie(w, userId)
	err = server.writeAutorizationHeader(w, userId)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
