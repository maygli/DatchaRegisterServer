package commandservice

import (
	"datcha/repository/channelrepository"
	"datcha/repository/commandrepository"
	"datcha/services/authservice/authinternalservice"
	"net/http"
)

type CommandService struct {
	cmdRepository     commandrepository.CommandRepositorier
	channelRepository channelrepository.ChannelRepositorier
	authService       *authinternalservice.AuthInternalService
}

func NewCommandService(authService *authinternalservice.AuthInternalService,
	cmdRepository commandrepository.CommandRepositorier,
	channelRepository channelrepository.ChannelRepositorier) *CommandService {
	service := CommandService{
		authService:       authService,
		cmdRepository:     cmdRepository,
		channelRepository: channelRepository,
	}
	return &service
}

func (service *CommandService) RegisterHandlers(r *http.ServeMux) {
	r.Handle("PUT /command", service.authService.JwtAuth(http.HandlerFunc(service.putCommandHandle)))
}
