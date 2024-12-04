package auththirdparty

import (
	"datcha/servercommon"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/wesleym/telegramwidget"
)

const (
	TELEGRAM_CONFIGURATION_PATH = "$.auth.telegram"
	TELEGRAM_REDIRECT_URL       = "/auth_pages/telegram_login.html"
	TELEGRAM_USER_KEY           = "id"
	TELEGRAM_FIRST_NAME_KEY     = "first_name"
	TELEGRAM_LAST_NAME_KEY      = "last_name"
	TELEGRAM_AVATAR_KEY         = "photo_url"
	TELEGRAM_HASH_KEY           = "hash"
)

type TelegramUser struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type TelegramAuthService struct {
	BotId     string `json:"bot_id" env:"${SERVER_NAME}_AUTH_TELEGRAM_BOT_ID"`
	BotSecret string `json:"bot_secret" env:"${SERVER_NAME}_AUTH_TELEGRAM_BOT_SECRET"`
}

func NewTelegramAuthService(cfgReader *servercommon.ConfigurationReader) (*TelegramAuthService, error) {
	service := TelegramAuthService{}
	err := cfgReader.ReadConfiguration(&service, TELEGRAM_CONFIGURATION_PATH)
	return &service, err
}

func (service TelegramAuthService) telegramAuthLoginHandle(w http.ResponseWriter, r *http.Request) {
	slog.Error("Telegram login")
	http.Redirect(w, r, TELEGRAM_REDIRECT_URL, http.StatusSeeOther)
}

func (service TelegramAuthService) telegramAuthCallbackHandle(w http.ResponseWriter, r *http.Request) {
	user, err := telegramwidget.ConvertAndVerifyForm(r.URL.Query(), service.BotSecret)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parse Telegram calback response.Error: %s", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	/*	user := TelegramUser{}
		user.UserId = r.URL.Query().Get(TELEGRAM_USER_KEY)
		if( user.UserId == "" ) {
			slog.Error("Telegram callback request doesn't contains user id")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		user.FirstName = r.URL.Query().Get(TELEGRAM_FIRST_NAME_KEY)
		user.LastName = r.URL.Query().Get(TELEGRAM_LAST_NAME_KEY)
		user.Avatar = r.URL.Query().Get(TELEGRAM_AVATAR_KEY)*/
	slog.Info(fmt.Sprintf("Login using telegram %+v", user))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (service TelegramAuthService) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/telegram/login", service.telegramAuthLoginHandle)
	mux.HandleFunc("GET /auth/telegram/callback", service.telegramAuthCallbackHandle)
}
