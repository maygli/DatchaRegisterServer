package authinternalservice

import (
	"context"
	"datcha/servercommon"
	"fmt"
	"log/slog"
	"net/http"
)

func (service *AuthInternalService) JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := service.getUserFromRequest(r)
		if err != nil {
			slog.Error(fmt.Sprintf("Attemt of non autorised access. Error: %s", err.Error()))
			http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
			return
		}
		slog.Info(fmt.Sprintf("Request with user=%s", user.Name))
		newContext := context.WithValue(r.Context(), servercommon.USER_CONTEXT_KEY, user)
		next.ServeHTTP(w, r.WithContext(newContext))
	})
}
