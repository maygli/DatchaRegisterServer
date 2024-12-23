package auththirdparty

import (
	"context"
	"datcha/servercommon"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

const (
	AUTH_YANDEX_SCOPE_EMAIL   = "login:email"
	AUTH_YANDEX_SCOPE_INFO    = "login:info"
	AUTH_YANDEX_SCOPE_AVATAR  = "login:avatar"
	AUTH_YANDEX_USER_INFO_URL = "https://login.yandex.ru/info?format=json"
	YANDEX_OAUTH              = "OAuth "
)

type YandexUser struct {
	UserId    string `json:"id"`
	Email     string `json:"default_email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"default_avatar_id"`
	Login     string `json:"login"`
}

type YandexAuthConfigure struct {
	YandexClientId     string `json:"yandex_client_id" env:"${SERVER_NAME}_AUTH_YANDEX_CLIENT_ID"`
	YandexClientSecret string `json:"yandex_client_secret" env:"${SERVER_NAME}_AUTH_YANDEX_CLIENT_SECRET"`
	ServerAddress      string `json:"yandex_address" env:"${SERVER_NAME}_SERVER_ADDRESS"`
	LoginEndPoint      string `json:"yandex_login_endpoint" env:"${SERVER_NAME}_YANDEX_LOGIN_ENDPOINT" default:"/auth/yandex/login"`
	CallbackEndPoint   string `json:"yandex_callback_endpoint" env:"${SERVER_NAME}_YANDEX_CALLBACK_ENDPOINT" default:"/auth/yandex/callback"`
}

type YandexAuthService struct {
	AuthThirdPartyBase
	yandexLoginConfig oauth2.Config
}

func NewYandexAuthService(config *YandexAuthConfigure, processor CompleteAuthProcessor) *YandexAuthService {
	service := YandexAuthService{
		AuthThirdPartyBase: AuthThirdPartyBase{
			LoginEndPoint:    config.LoginEndPoint,
			CallbackEndPoint: config.CallbackEndPoint,
			AuthProcessor:    processor,
		},
	}
	service.yandexLoginConfig = oauth2.Config{
		RedirectURL:  config.ServerAddress + service.GetCallbackEndpoint(service.GetServiceName()),
		ClientID:     config.YandexClientId,
		ClientSecret: config.YandexClientSecret,
		Scopes:       []string{AUTH_YANDEX_SCOPE_INFO, AUTH_YANDEX_SCOPE_EMAIL, AUTH_YANDEX_SCOPE_AVATAR},
		Endpoint:     yandex.Endpoint,
	}
	return &service
}

func (service YandexAuthService) GetServiceName() string {
	return servercommon.YANDEX_SERVICE_NAME
}

func (service YandexAuthService) yandexAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	ProcessLogin(w, r, service.yandexLoginConfig)
}

func (service YandexAuthService) getUser(token *oauth2.Token) (YandexUser, error) {
	client := service.yandexLoginConfig.Client(context.Background(), token)
	req, err := http.NewRequest(http.MethodGet, AUTH_YANDEX_USER_INFO_URL, nil)
	if err != nil {
		return YandexUser{}, err
	}
	req.Header.Set(AUTH_HEADER, YANDEX_OAUTH+token.AccessToken)
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

func YandexToAuthUser(user YandexUser) AuthThirdPartyUser {
	resUser := AuthThirdPartyUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		UserId:    user.UserId,
		Service:   servercommon.YANDEX_SERVICE_NAME,
	}
	return resUser
}

func (service YandexAuthService) yandexAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	slog.Info("Yandex callback called")
	token, err := GetToken(r, service.yandexLoginConfig)
	if err != nil {
		slog.Error(fmt.Sprintf("Error get token. Error: %s", err.Error()))
		service.ProcessError(err, w, r)
		return
	}
	slog.Info(fmt.Sprintf("Get yandex token %+v", token))
	user, err := service.getUser(token)
	if err != nil {
		slog.Error(fmt.Sprintf("can't get yandex user data. Error: %s", err.Error()))
		service.ProcessError(err, w, r)
		return
	}
	slog.Info(fmt.Sprintf("User data %+v", user))
	yandexUser := YandexToAuthUser(user)
	err = service.ProcessSuccess(yandexUser, w, r)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (service YandexAuthService) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET "+service.GetLoginEndpoint(service.GetServiceName()), service.yandexAuthLoginHandle)
	mux.HandleFunc("GET "+service.GetCallbackEndpoint(service.GetServiceName()), service.yandexAuthCallbackHandle)
}
