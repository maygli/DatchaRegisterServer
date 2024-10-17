package authentification

import (
	"context"
	"datcha/servercommon"
	"log"
	"net/http"
)

func (server *AuthServer) JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := server.getUserFromRequest(r)
		if err != nil {
			log.Println("Attemt of non autorised access. Error: " + err.Error())
			http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
			return
		}
		log.Println("Request with user=" + user.Name)
		newContext := context.WithValue(r.Context(), servercommon.USER_CONTEXT_KEY, user)
		next.ServeHTTP(w, r.WithContext(newContext))
	})
}
