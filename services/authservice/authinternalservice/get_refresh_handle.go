package authinternalservice

import (
	"datcha/servercommon"
	"net/http"
)

func (server *AuthInternalService) refreshGetHandle(w http.ResponseWriter, r *http.Request) {
	user, err := server.getUserFromRequest(r)
	if err != nil {
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	err = server.WriteAutorizationHeader(w, user.ID)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	err = server.writeStatusResponse(w, user.AccStatus)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
