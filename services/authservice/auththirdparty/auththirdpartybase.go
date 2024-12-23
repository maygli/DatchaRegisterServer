package auththirdparty

import (
	"fmt"
	"net/http"
)

const (
	LOGIN_END_POINT_TEMPLATE    = "/auth/%s/login"
	CALLBACK_END_POINT_TEMPLATE = "/auth/%s/callback"
)

type AuthThirdPartyUser struct {
	FirstName string
	LastName  string
	Email     string
	Avatar    string
	UserId    string
	Service   string
}

type CompleteAuthProcessor interface {
	ProcessSuceess(user any, w http.ResponseWriter, r *http.Request) error
	ProcessError(err error, w http.ResponseWriter, r *http.Request) error
}

type IAuthThirdPartyService interface {
	RegisterHandlers(mux *http.ServeMux)
	GetServiceName() string
}

type AuthThirdPartyBase struct {
	LoginEndPoint    string
	CallbackEndPoint string
	AuthProcessor    CompleteAuthProcessor
}

func (user AuthThirdPartyUser) GetName() string {
	name := user.FirstName
	if name != "" && user.LastName != "" {
		name += " "
	}
	name += user.LastName
	return name
}

func (service AuthThirdPartyBase) GetLoginEndpoint(serviceName string) string {
	if service.LoginEndPoint != "" {
		return service.LoginEndPoint
	}
	return fmt.Sprintf(LOGIN_END_POINT_TEMPLATE, serviceName)
}

func (service AuthThirdPartyBase) GetCallbackEndpoint(serviceName string) string {
	if service.CallbackEndPoint != "" {
		return service.CallbackEndPoint
	}
	return fmt.Sprintf(CALLBACK_END_POINT_TEMPLATE, serviceName)
}

func (service AuthThirdPartyBase) ProcessError(err error, w http.ResponseWriter, r *http.Request) {
	if service.AuthProcessor != nil {
		service.AuthProcessor.ProcessError(err, w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (service AuthThirdPartyBase) ProcessSuccess(user AuthThirdPartyUser, w http.ResponseWriter, r *http.Request) error {
	var err error = nil
	if service.AuthProcessor != nil {
		err := service.AuthProcessor.ProcessSuceess(user, w, r)
		if err == nil {
			return err
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return err
}
