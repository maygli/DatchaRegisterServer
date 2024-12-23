package profileservice

import (
	"datcha/repository/authrepository"
	"datcha/repository/profilerepository"
	"datcha/services/authservice/authinternalservice"
	"net/http"
)

type ProfileData struct {
	ID     uint   `json:"profile_id" schema:"profile_id"`
	Name   string `json:"user_name" schema:"user_name"`
	Avatar string `json:"avatar" schema:"avatar"`
}

type ProfileService struct {
	userRepository    authrepository.AuthRepositorier
	profileRepository profilerepository.ProfileRepositorier
	authService       *authinternalservice.AuthInternalService
}

func NewProfileService(authRep authrepository.AuthRepositorier,
	profileRep profilerepository.ProfileRepositorier,
	authService *authinternalservice.AuthInternalService) *ProfileService {
	service := ProfileService{
		userRepository:    authRep,
		profileRepository: profileRep,
		authService:       authService,
	}
	return &service
}

func (profileData ProfileData) isValid() bool {
	return profileData.Name != ""
}

func (service *ProfileService) RegisterHandlers(r *http.ServeMux) {
	r.Handle("GET /profile", service.authService.JwtAuth(http.HandlerFunc(service.getProfileHandle)))
	r.Handle("PUT /profile", service.authService.JwtAuth(http.HandlerFunc(service.putProfileHandle)))
}
