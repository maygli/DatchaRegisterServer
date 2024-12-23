package authinternalservice

import (
	"datcha/datamodel"
	"datcha/servercommon"

	"log"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func (server *AuthInternalService) registerPostHandle(w http.ResponseWriter, r *http.Request) {
	userData, err := server.getUserData(r)
	if err != nil {
		log.Println("Error. Can't parse user data. Error=", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = server.validateUserName(userData.UserName)
	if err != nil {
		log.Println("Error. Validation user data failed. Error=", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(userData.Password) < MIN_PASSWORD_LENGHT {
		log.Println("Error. Password is too short.")
		http.Error(w, servercommon.ERROR_PASSWORD_TOO_SHORT, http.StatusBadRequest)
		return
	}
	_, err = mail.ParseAddress(userData.Email)
	if err != nil {
		log.Println("Error. Email address is invalid. Error=" + err.Error())
		http.Error(w, servercommon.ERROR_EMAIL_INVALID, http.StatusBadRequest)
		return
	}
	genHashPwd, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error. Can't generate password hash. Error=" + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	accStatus := datamodel.AS_NEED_CONFIRM
	if !server.shouldSendEmailConfirmation() {
		accStatus = datamodel.AS_CONFIRMED
	}
	user := datamodel.User{
		Name:      userData.UserName,
		Email:     userData.Email,
		Password:  string(genHashPwd),
		AccStatus: accStatus,
	}
	err = server.repository.AddUser(&user)
	if err != nil {
		log.Println("Error. Can't add user. Error=" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = server.WriteAutorizationHeader(w, user.ID)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	server.SendAuthCoockie(w, user.ID)
	if server.shouldSendEmailConfirmation() {
		err := server.SendConfirmationEmail(user.ID, user.Name, user.Email, userData.Locale)
		if err != nil {
			log.Println("Error. Can't send confirmation email. Error=" + err.Error())
			err = server.repository.DeleteUser(user.ID)
			if err != nil {
				log.Println("Can't delete user. Error: " + err.Error())
			}
			http.Error(w, servercommon.ERROR_EMAIL_INVALID, http.StatusBadRequest)
		}
	}
	err = server.writeStatusResponse(w, user.AccStatus)
	if err != nil {
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
