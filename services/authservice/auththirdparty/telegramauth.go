package auththirdparty

import (
	"datcha/servercommon"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/wesleym/telegramwidget"
)

const (
	TELEGRAM_USER_KEY       = "id"
	TELEGRAM_FIRST_NAME_KEY = "first_name"
	TELEGRAM_LAST_NAME_KEY  = "last_name"
	TELEGRAM_AVATAR_KEY     = "photo_url"
	TELEGRAM_HASH_KEY       = "hash"
)

type TelegramUser struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type TelegramAuthConfigure struct {
	BotId            string `json:"bot_id" env:"${SERVER_NAME}_AUTH_TELEGRAM_BOT_ID"`
	BotSecret        string `json:"bot_secret" env:"${SERVER_NAME}_AUTH_TELEGRAM_BOT_SECRET"`
	ServerAddress    string `json:"server_address" env:"${SERVER_NAME}_SERVER_ADDRESS"`
	LoginEndPoint    string `json:"telegram_login_endpoint" env:"${SERVER_NAME}_TELEGRAM_LOGIN_ENDPOINT" default:"/auth/telegram/login"`
	CallbackEndPoint string `json:"telegram_callback_endpoint" env:"${SERVER_NAME}_TELEGRAM_CALLBACK_ENDPOINT" default:"/telegram/google/callback"`
	LoginPageUrl     string `json:"telegram_login_page_url" env:"${SERVER_NAME}_TELEGRAM_LOGIN_PAGE_URL" default:"/auth_pages/telegram-login.html"`
}

type TelegramAuthService struct {
	AuthThirdPartyBase
	BotId        string
	BotSecret    string
	LoginPageUrl string
}

func NewTelegramAuthService(config *TelegramAuthConfigure, processor CompleteAuthProcessor) *TelegramAuthService {
	service := TelegramAuthService{
		AuthThirdPartyBase: AuthThirdPartyBase{
			LoginEndPoint:    config.LoginEndPoint,
			CallbackEndPoint: config.CallbackEndPoint,
			AuthProcessor:    processor,
		},
		BotId:        config.BotId,
		BotSecret:    config.BotSecret,
		LoginPageUrl: config.LoginPageUrl,
	}
	return &service
}

func (service TelegramAuthService) GetServiceName() string {
	return servercommon.TELEGRAM_SERVICE_NAME
}

func (service TelegramAuthService) telegramAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	slog.Error("Telegram login")
	http.Redirect(w, r, service.LoginPageUrl, http.StatusSeeOther)
}

func telegramToAuthUser(user telegramwidget.User) AuthThirdPartyUser {
	resUser := AuthThirdPartyUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.PhotoURL.Host,
		UserId:    strconv.FormatInt(user.ID, 10),
		Service:   servercommon.TELEGRAM_SERVICE_NAME,
	}
	return resUser
}

func (service TelegramAuthService) telegramAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	user, err := telegramwidget.ConvertAndVerifyForm(r.URL.Query(), service.BotSecret)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parse Telegram calback response.Error: %s", err.Error()))
		service.ProcessError(err, w, r)
		return
	}
	slog.Info(fmt.Sprintf("Login using telegram %+v", user))
	telegramUser := telegramToAuthUser(user)
	err = service.ProcessSuccess(telegramUser, w, r)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (service TelegramAuthService) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET "+service.GetLoginEndpoint(service.GetServiceName()), service.telegramAuthLoginHandle)
	mux.HandleFunc("GET "+service.GetCallbackEndpoint(service.GetCallbackEndpoint(service.GetServiceName())), service.telegramAuthCallbackHandle)
}
