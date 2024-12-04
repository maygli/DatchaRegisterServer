package auththirdparty

import (
	"context"
	"datcha/servercommon"
	"datcha/services/authservice/authcommon"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	AUTH_GOOGLE_CONFIGURATION_PATH = "$.auth.google"
	AUTH_GOOGLE_CALLBACK_ENDPOINT  = "/auth/google/callback"
	AUTH_GOOGLE_EMAIL_SCOPE        = "https://www.googleapis.com/auth/userinfo.email"
	AUTH_GOOGLE_PROFILE_SCOPE      = "https://www.googleapis.com/auth/userinfo.profile"
	GOOGLE_GET_USER_METHOD         = http.MethodGet
	GOOGLE_USER_PROFILE_DATA_URL   = "https://www.googleapis.com/oauth2/v2/userinfo"
	GOOGLE_OAUTH_API_URL           = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type GoogleUser struct {
	UserId    string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Avatar    string `json:"picture"`
}

type GoogleAuthService struct {
	GoogleClientId     string `json:"google_client_id" env:"${SERVER_NAME}_AUTH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `json:"google_client_secret" env:"${SERVER_NAME}_AUTH_GOOGLE_CLIENT_SECRET"`
	ServerAddress      string `json:"server_address" env:"${SERVER_NAME}_SERVER_ADDRESS"`
	googleLoginConfig  oauth2.Config
}

func NewGoogleAuthService(cfgReader *servercommon.ConfigurationReader) (*GoogleAuthService, error) {
	service := GoogleAuthService{}
	err := cfgReader.ReadConfiguration(&service, AUTH_GOOGLE_CONFIGURATION_PATH)
	if err != nil {
		return &service, err
	}
	service.googleLoginConfig = oauth2.Config{
		RedirectURL:  service.ServerAddress + AUTH_GOOGLE_CALLBACK_ENDPOINT,
		ClientID:     service.GoogleClientId,
		ClientSecret: service.GoogleClientSecret,
		Scopes: []string{AUTH_GOOGLE_EMAIL_SCOPE,
			AUTH_GOOGLE_PROFILE_SCOPE},
		Endpoint: google.Endpoint,
	}
	return &service, err
}

func (service GoogleAuthService) googleAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	authcommon.ProcessLogin(w, r, service.googleLoginConfig)
}

func (service GoogleAuthService) getUserDataFromGoogle(token *oauth2.Token) (GoogleUser, error) {
	client := service.googleLoginConfig.Client(context.Background(), token)
	response, err := client.Get(GOOGLE_OAUTH_API_URL + token.AccessToken)
	if err != nil {
		return GoogleUser{}, err
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return GoogleUser{}, err
	}
	user := GoogleUser{}
	slog.Info(fmt.Sprintf("Request content=%s", string(content)))
	err = json.Unmarshal(content, &user)
	return user, err
}

func (service GoogleAuthService) googleAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	token, err := authcommon.GetToken(r, service.googleLoginConfig)
	if err != nil {
		slog.Error(fmt.Sprintf("Error get token. Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user, err := service.getUserDataFromGoogle(token)
	if err != nil {
		slog.Error("Error get token")
		return
	}
	slog.Info(fmt.Sprintf("user=%+v", user))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (service GoogleAuthService) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/google/login", service.googleAuthLoginHandle)
	mux.HandleFunc("GET "+AUTH_GOOGLE_CALLBACK_ENDPOINT, service.googleAuthCallbackHandle)
}
