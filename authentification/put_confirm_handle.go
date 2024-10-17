package authentification

import (
	"datcha/servercommon"
	"log"
	"net/http"
)

func (server *AuthServer) confirmPutHandle(w http.ResponseWriter, r *http.Request) {
	if !server.shouldSendEmailConfirmation() {
		log.Println("Server doesn't support email confirmation but got request")
		// Or should we send error?
		//		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	user, err := server.getUserFromRequest(r)
	if err != nil {
		log.Println("Error. Can't get user data from request. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
	}
	locale := "ru"
	err = server.SendConfirmationEmail(user.ID, user.Name, user.Email, locale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
