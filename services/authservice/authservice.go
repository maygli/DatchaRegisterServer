package authservice

import (
	"datcha/repository/authrepository"
	"datcha/repository/profilerepository"
	"datcha/servercommon"
	"datcha/services/authservice/authinternalservice"
	"datcha/services/authservice/auththirdparty"
	"net/http"
)

const (
	AUTH_CONFIGURATION_PATH = "$.auth"
)

type AuthService struct {
	ThirdPartyServicesList []string `json:"thirdparty_services" env:"${SERVER_NAME}_THIRD_PARTY_SERVICES"`
	ThirdPartyServices     *auththirdparty.AuthThirdPartyService
	InternalService        *authinternalservice.AuthInternalService
}

func NewAuthService(cfgReader *servercommon.ConfigurationReader,
	authRepository authrepository.AuthRepositorier,
	profileRepository profilerepository.ProfileRepositorier) (*AuthService, error) {
	service := AuthService{}
	err := cfgReader.ReadConfiguration(&service, AUTH_CONFIGURATION_PATH)
	if err != nil {
		return nil, err
	}
	service.InternalService, err = authinternalservice.NewAuthServer(cfgReader, authRepository)
	if err != nil {
		return nil, err
	}
	if len(service.ThirdPartyServicesList) > 0 {
		service.ThirdPartyServices, err = auththirdparty.NewAuthThirdPartyService(cfgReader, service.ThirdPartyServicesList,
			authRepository, profileRepository, service.InternalService)
		if err != nil {
			return nil, err
		}
	}
	return &service, err
}

func (service AuthService) RegisterHandlers(mux *http.ServeMux) {
	if service.ThirdPartyServices != nil {
		service.ThirdPartyServices.RegisterHandlers(mux)
	}
	if service.InternalService != nil {
		service.InternalService.RegisterHandlers(mux)
	}
}
