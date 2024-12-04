package channelservice

import (
	"datcha/repository/channelrepository"
	"datcha/servercommon"
	"datcha/services/authservice/authinternalservice"
	"net/http"
)

type ChannelService struct {
	channelRepository channelrepository.ChannelRepositorier
	authService       *authinternalservice.AuthInternalService
}

func NewChannelService(authService *authinternalservice.AuthInternalService,
	rep channelrepository.ChannelRepositorier) *ChannelService {
	service := ChannelService{
		channelRepository: rep,
		authService:       authService,
	}
	return &service
}

func (service *ChannelService) RegisterHandlers(r *http.ServeMux) {
	r.Handle("GET /channel/data/{"+servercommon.CHANNEL_ID_KEY+"}", service.authService.JwtAuth(http.HandlerFunc(service.getChannelDataHandle)))
}
