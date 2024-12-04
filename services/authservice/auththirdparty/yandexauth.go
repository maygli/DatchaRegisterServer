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
	"golang.org/x/oauth2/yandex"
)

const (
	AUTH_YANDEX_CONFIGURATION_PATH = "$.auth.yandex"
	AUTH_YANDEX_CALLBACK_ENDPOINT  = "/auth/yandex/callback"
	AUTH_YANDEX_SCOPE_EMAIL        = "login:email"
	AUTH_YANDEX_SCOPE_INFO         = "login:info"
	AUTH_YANDEX_SCOPE_AVATAR       = "login:avatar"
	AUTH_YANDEX_USER_INFO_URL      = "https://login.yandex.ru/info?format=json"
	YANDEX_OAUTH                   = "OAuth "
)

type YandexUser struct {
	UserId    string `json:"id"`
	Emails    string `json:"default_email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"default_avatar_id"`
	Login     string `json:"login"`
}

type YandexAuthService struct {
	YandexClientId     string `json:"yandex_client_id" env:"${SERVER_NAME}_AUTH_YANDEX_CLIENT_ID"`
	YandexClientSecret string `json:"yandex_client_secret" env:"${SERVER_NAME}_AUTH_YANDEX_CLIENT_SECRET"`
	ServerAddress      string `json:"server_address" env:"${SERVER_NAME}_SERVER_ADDRESS"`
	yandexLoginConfig  oauth2.Config
}

func NewYandexAuthService(cfgReader *servercommon.ConfigurationReader) (*YandexAuthService, error) {
	service := YandexAuthService{}
	err := cfgReader.ReadConfiguration(&service, AUTH_YANDEX_CONFIGURATION_PATH)
	if err != nil {
		return nil, err
	}
	service.yandexLoginConfig = oauth2.Config{
		RedirectURL:  service.ServerAddress + AUTH_YANDEX_CALLBACK_ENDPOINT,
		ClientID:     service.YandexClientId,
		ClientSecret: service.YandexClientSecret,
		Scopes:       []string{AUTH_YANDEX_SCOPE_INFO, AUTH_YANDEX_SCOPE_EMAIL, AUTH_YANDEX_SCOPE_AVATAR},
		Endpoint:     yandex.Endpoint,
	}
	return &service, err
}

func (service YandexAuthService) yandexAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	authcommon.ProcessLogin(w, r, service.yandexLoginConfig)
}

func (service YandexAuthService) getUser(token *oauth2.Token) (YandexUser, error) {
	client := service.yandexLoginConfig.Client(context.Background(), token)
	req, err := http.NewRequest(http.MethodGet, AUTH_YANDEX_USER_INFO_URL, nil)
	if err != nil {
		return YandexUser{}, err
	}
	req.Header.Set(authcommon.AUTH_HEADER, YANDEX_OAUTH+token.AccessToken)
	response, err := client.Do(req)
	if err != nil {
		return YandexUser{}, err
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return YandexUser{}, err
	}
	user := YandexUser{}
	slog.Info(fmt.Sprintf("Request content=%s", string(content)))
	err = json.Unmarshal(content, &user)
	return user, err
}

func (service YandexAuthService) yandexAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	slog.Info("Yandex callback called")
	token, err := authcommon.GetToken(r, service.yandexLoginConfig)
	if err != nil {
		slog.Error(fmt.Sprintf("Error get token. Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	slog.Info(fmt.Sprintf("Get yandex token %+v", token))
	user, err := service.getUser(token)
	if err != nil {
		slog.Error(fmt.Sprintf("can't get yandex user data. Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	slog.Info(fmt.Sprintf("User data %+v", user))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (service YandexAuthService) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/yandex/login", service.yandexAuthLoginHandle)
	mux.HandleFunc("GET "+AUTH_YANDEX_CALLBACK_ENDPOINT, service.yandexAuthCallbackHandle)
}
