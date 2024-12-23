package authinternalservice

import (
	"datcha/servercommon"
	"log"
	"net/http"
	"strings"
)

func (server *AuthInternalService) refreshPostHandle(w http.ResponseWriter, r *http.Request) {
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
	server.SendAuthCoockie(w, userId)
	err = server.WriteAutorizationHeader(w, userId)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
