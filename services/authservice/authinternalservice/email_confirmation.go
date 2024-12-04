package authinternalservice

import (
	"bytes"
	"datcha/servercommon"
	"errors"
	"log"
	"net/smtp"
	"os"
	"path"
	"strconv"
	"text/template"
)

const (
	CONFIRMATION_EMAIL_TEMPLATE = "confirm"
	TEMPLATE_EXTENSION          = ".templ"
)

func (server *AuthInternalService) shouldSendEmailConfirmation() bool {
	if server.SMTPServerURL != "" {
		return true
	}
	return false
}

func (server *AuthInternalService) SendConfirmationEmail(userId uint, userName string, email string, locale string) error {
	if server.SMTPServerURL == "" {
		return nil
	}

	// Receiver email address.
	to := []string{
		email,
	}

	// Message.
	message, err := server.generateConfirmationMessage(userId, userName, locale)
	if err != nil {
		log.Println("Error. Can't generate confirmation message. Error: " + err.Error())
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	// Authentication.
	auth := smtp.PlainAuth("", server.ConfiramationEmail, server.EmailPassword, server.SMTPServerURL)

	// Sending email.
	err = smtp.SendMail(server.SMTPServerURL+":"+strconv.Itoa(server.SMTPServerPort), auth, server.ConfiramationEmail, to, message.Bytes())
	if err != nil {
		log.Println("Error. Can't send confirmation mail. Error: ", err)
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	return nil
}

func (server *AuthInternalService) getMailTemplate(templFileName string, locale string) (string, error) {
	filePath := path.Join(server.EmailTemplatesFolder, templFileName+"_"+locale+TEMPLATE_EXTENSION)
	_, err := os.Stat(filePath)
	if err == nil {
		return filePath, nil
	}
	filePath = path.Join(server.EmailTemplatesFolder, templFileName+TEMPLATE_EXTENSION)
	_, err = os.Stat(filePath)
	return filePath, err
}

func (server *AuthInternalService) generateConfirmationMessage(userId uint, userName string, locale string) (bytes.Buffer, error) {
	var message bytes.Buffer
	templFileName, err := server.getMailTemplate(CONFIRMATION_EMAIL_TEMPLATE, locale)
	if err != nil {
		return message, err
	}
	templ, err := template.ParseFiles(templFileName)
	if err != nil {
		return message, err
	}
	token, err := server.generateToken(userId, server.ConfirmTokenAge, server.ConfirmSecretKey, CONFIRM_SUBJECT)
	err = templ.Execute(&message, struct {
		Name    string
		Product string
		Link    string
	}{
		Name:    userName,
		Product: "CrazyDatcha",
		Link:    "http://localhost:6080/confirm/" + token,
	})
	return message, err
}
