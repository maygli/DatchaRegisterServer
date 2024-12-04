package authinternalservice

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"log"
	"net/http"
)

func (server *AuthInternalService) confirmGetHandle(w http.ResponseWriter, r *http.Request) {
	confirmToken := r.PathValue(servercommon.CONFIRM_TOKEN_KEY)
	if confirmToken == "" {
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		log.Println("Error. Confirm request doesn't contains token")
		return
	}
	userId, err := server.getUserIdFromToken(confirmToken, CONFIRM_SUBJECT, server.ConfirmSecretKey)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user, err := server.repository.FindUser(userId)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err = server.repository.UpdateAccountStatus(user.ID, datamodel.AS_CONFIRMED)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("confirm user", user)
}
