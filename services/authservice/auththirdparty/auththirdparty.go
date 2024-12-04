package auththirdparty

import (
	"datcha/servercommon"
	"errors"
	"fmt"
	"net/http"
	"slices"
)

const (
	GOOGLE_SERVICE_NAME   = "google"
	VK_SERVICE_NAME       = "vk"
	TELEGRAM_SERVICE_NAME = "telegram"
	YANDEX_SERVICE_NAME   = "yandex"
)

type AuthThirdPartyService struct {
	AuthGoogleService   *GoogleAuthService
	AuthVkService       *VkAuthService
	AuthTelegramService *TelegramAuthService
	AuthYandexService   *YandexAuthService
}

func NewAuthThirdPartyService(cfgReader *servercommon.ConfigurationReader, servicesList []string) (*AuthThirdPartyService, error) {
	if len(servicesList) == 0 {
		return nil, nil
	}
	service := AuthThirdPartyService{}
	var err error
	if slices.Contains(servicesList, GOOGLE_SERVICE_NAME) {
		service.AuthGoogleService, err = NewGoogleAuthService(cfgReader)
		if err != nil {
			msg := fmt.Sprintf("Can't create Google authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
	}
	if slices.Contains(servicesList, VK_SERVICE_NAME) {
		service.AuthVkService, err = NewVkAuthService(cfgReader)
		if err != nil {
			msg := fmt.Sprintf("Can't create Vk authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
	}
	if slices.Contains(servicesList, TELEGRAM_SERVICE_NAME) {
		service.AuthTelegramService, err = NewTelegramAuthService(cfgReader)
		if err != nil {
			msg := fmt.Sprintf("Can't create Vk authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
	}
	if slices.Contains(servicesList, TELEGRAM_SERVICE_NAME) {
		service.AuthYandexService, err = NewYandexAuthService(cfgReader)
		if err != nil {
			msg := fmt.Sprintf("Can't create Yandex authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
	}
	return &service, nil
}

func (service AuthThirdPartyService) RegisterHandlers(mux *http.ServeMux) {
	if service.AuthGoogleService != nil {
		service.AuthGoogleService.RegisterHandlers(mux)
	}
	if service.AuthVkService != nil {
		service.AuthVkService.RegisterHandlers(mux)
	}
	if service.AuthTelegramService != nil {
		service.AuthTelegramService.RegisterHandlers(mux)
	}
	if service.AuthYandexService != nil {
		service.AuthYandexService.RegisterHandlers(mux)
	}
}
