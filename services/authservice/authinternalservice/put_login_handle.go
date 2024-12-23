package authinternalservice

import (
	"datcha/servercommon"
	"fmt"

	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (server *AuthInternalService) loginPutHandle(w http.ResponseWriter, r *http.Request) {
	userData, err := server.getUserData(r)
	if err != nil {
		log.Println("error parse login data. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusBadRequest)
		return
	}
	if userData.UserName == "" || userData.Password == "" {
		log.Println("username or(and) password is empty")
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusBadRequest)
		return
	}
	fmt.Println("user name=", userData.UserName)
	fmt.Println("password=", userData.Password)
	user, err := server.repository.FindUserByNameEmail(userData.UserName)
	if err != nil {
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	err = server.WriteAutorizationHeader(w, user.ID)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	server.SendAuthCoockie(w, user.ID)
	err = server.writeStatusResponse(w, user.AccStatus)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	//	http.Redirect(w, r, "/", http.StatusSeeOther)
}
