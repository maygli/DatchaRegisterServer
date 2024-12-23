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
	"golang.org/x/oauth2/google"
)

const (
	AUTH_GOOGLE_EMAIL_SCOPE      = "https://www.googleapis.com/auth/userinfo.email"
	AUTH_GOOGLE_PROFILE_SCOPE    = "https://www.googleapis.com/auth/userinfo.profile"
	GOOGLE_GET_USER_METHOD       = http.MethodGet
	GOOGLE_USER_PROFILE_DATA_URL = "https://www.googleapis.com/oauth2/v2/userinfo"
	GOOGLE_OAUTH_API_URL         = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type GoogleUser struct {
	UserId    string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Avatar    string `json:"picture"`
}

type GoogleAuthConfigure struct {
	GoogleClientId     string `json:"google_client_id" env:"${SERVER_NAME}_AUTH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `json:"google_client_secret" env:"${SERVER_NAME}_AUTH_GOOGLE_CLIENT_SECRET"`
	ServerAddress      string `json:"server_address" env:"${SERVER_NAME}_SERVER_ADDRESS"`
	LoginEndPoint      string `json:"google_login_endpoint" env:"${SERVER_NAME}_GOOGLE_LOGIN_ENDPOINT" default:"/auth/google/login"`
	CallbackEndPoint   string `json:"google_callback_endpoint" env:"${SERVER_NAME}_GOOGLE_CALLBACK_ENDPOINT" default:"/auth/google/callback"`
}

type GoogleAuthService struct {
	AuthThirdPartyBase
	googleLoginConfig oauth2.Config
}

func NewGoogleAuthService(config *GoogleAuthConfigure, processor CompleteAuthProcessor) *GoogleAuthService {
	service := GoogleAuthService{
		AuthThirdPartyBase: AuthThirdPartyBase{
			LoginEndPoint:    config.LoginEndPoint,
			CallbackEndPoint: config.CallbackEndPoint,
			AuthProcessor:    processor,
		},
	}
	service.googleLoginConfig = oauth2.Config{
		RedirectURL:  config.ServerAddress + service.GetCallbackEndpoint(service.GetServiceName()),
		ClientID:     config.GoogleClientId,
		ClientSecret: config.GoogleClientSecret,
		Scopes: []string{AUTH_GOOGLE_EMAIL_SCOPE,
			AUTH_GOOGLE_PROFILE_SCOPE},
		Endpoint: google.Endpoint,
	}
	return &service
}

func (service GoogleAuthService) GetServiceName() string {
	return servercommon.GOOGLE_SERVICE_NAME
}

func (service GoogleAuthService) googleAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	ProcessLogin(w, r, service.googleLoginConfig)
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

func googleToAuthUser(user GoogleUser) AuthThirdPartyUser {
	resUser := AuthThirdPartyUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		UserId:    user.UserId,
		Service:   servercommon.GOOGLE_SERVICE_NAME,
	}
	return resUser
}

func (service GoogleAuthService) googleAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	token, err := GetToken(r, service.googleLoginConfig)
	if err != nil {
		slog.Error(fmt.Sprintf("Error get token. Error: %s", err.Error()))
		service.ProcessError(err, w, r)
		return
	}
	user, err := service.getUserDataFromGoogle(token)
	if err != nil {
		slog.Error("Error get token. Error: " + err.Error())
		service.ProcessError(err, w, r)
		return
	}
	slog.Info(fmt.Sprintf("user=%+v", user))
	gUser := googleToAuthUser(user)
	err = service.ProcessSuccess(gUser, w, r)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (service GoogleAuthService) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET "+service.GetLoginEndpoint(service.GetServiceName()), service.googleAuthLoginHandle)
	mux.HandleFunc("GET "+service.GetCallbackEndpoint(service.GetServiceName()), service.googleAuthCallbackHandle)
}
