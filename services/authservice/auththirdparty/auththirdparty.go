package auththirdparty

import (
	"datcha/repository/authrepository"
	"datcha/repository/profilerepository"
	"datcha/servercommon"
	"datcha/services/authservice/authinternalservice"
	"errors"
	"fmt"
	"net/http"
	"slices"
)

const (
	GOOGLE_CONFIGURATION_PATH   = "$.auth.google"
	VK_CONFIGURATION_PATH       = "$.auth.vk"
	TELEGRAM_CONFIGURATION_PATH = "$.auth.telegram"
	YANDEX_CONFIGURATION_PATH   = "$.auth.yandex"
)

type AuthThirdPartyService struct {
	AuthGoogleService   *GoogleAuthService
	AuthVkService       *VkAuthService
	AuthTelegramService *TelegramAuthService
	AuthYandexService   *YandexAuthService
	authRepository      authrepository.AuthRepositorier
	profileRepository   profilerepository.ProfileRepositorier
	intServive          *authinternalservice.AuthInternalService
}

type ThirdPartyAuthService struct {
}

func NewAuthThirdPartyService(cfgReader *servercommon.ConfigurationReader, servicesList []string,
	authRep authrepository.AuthRepositorier, profileRep profilerepository.ProfileRepositorier,
	intService *authinternalservice.AuthInternalService) (*AuthThirdPartyService, error) {
	if len(servicesList) == 0 {
		return nil, nil
	}
	service := AuthThirdPartyService{
		authRepository:    authRep,
		profileRepository: profileRep,
		intServive:        intService,
	}
	var err error
	if slices.Contains(servicesList, servercommon.GOOGLE_SERVICE_NAME) {
		googleConfig := GoogleAuthConfigure{}
		err = cfgReader.ReadConfiguration(&googleConfig, GOOGLE_CONFIGURATION_PATH)
		if err != nil {
			msg := fmt.Sprintf("Can't create Google authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
		service.AuthGoogleService = NewGoogleAuthService(&googleConfig, &service)
	}
	if slices.Contains(servicesList, servercommon.VK_SERVICE_NAME) {
		vkConfig := VkAuthConfigure{}
		err = cfgReader.ReadConfiguration(&vkConfig, VK_CONFIGURATION_PATH)
		if err != nil {
			msg := fmt.Sprintf("Can't create Vk authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
		service.AuthVkService = NewVkAuthService(&vkConfig, &service)
	}
	if slices.Contains(servicesList, servercommon.TELEGRAM_SERVICE_NAME) {
		telegramConfig := TelegramAuthConfigure{}
		err = cfgReader.ReadConfiguration(&telegramConfig, TELEGRAM_CONFIGURATION_PATH)
		if err != nil {
			msg := fmt.Sprintf("Can't create Telegram authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
		service.AuthTelegramService = NewTelegramAuthService(&telegramConfig, &service)
	}
	if slices.Contains(servicesList, servercommon.TELEGRAM_SERVICE_NAME) {
		yandexConfig := YandexAuthConfigure{}
		err = cfgReader.ReadConfiguration(&yandexConfig, YANDEX_CONFIGURATION_PATH)
		if err != nil {
			msg := fmt.Sprintf("Can't create Yandex authentification service. Error: %s", err.Error())
			return nil, errors.New(msg)
		}
		service.AuthYandexService = NewYandexAuthService(&yandexConfig, &service)
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
