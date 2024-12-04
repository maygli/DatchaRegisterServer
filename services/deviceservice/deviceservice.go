package deviceservice

import (
	"datcha/repository/channelrepository"
	"datcha/repository/commandrepository"
	"datcha/repository/devicerepository"
	"datcha/servercommon"
	"net/http"
	"sync"
)

const (
	DEVICE_SERVICE_CONFIG_PATH string = "$.deviceserver"
)

type DeviceService struct {
	MapMutex          sync.Mutex
	deviceRepository  devicerepository.DeviceRepositorier
	channelRepository channelrepository.ChannelRepositorier
	commandRepository commandrepository.CommandRepositorier
	Issuer            string `json:"issuer" env:"${SERVER_NAME}_DEVICE_ISSUER" default:"crazydatchahome"`
	StateSubject      string `json:"state_subject" env:"${SERVER_NAME}_DEVICE_STATE_SUBJECT" default:"devicesstate"`
	StateSecretKey    string `json:"state_secret" env:"${SERVER_NAME}_DEVICE_STATE_SECRET" default:",g1XPz){ubdSK8{_(%rra[&t[rWSpEVV"`
	DeviceUser        string `json:"device_user" env:"${SERVER_NAME}_DEVICE_USER" default:"device"`
	DevicePassword    string `json:"device_password" env:"${SERVER_NAME}_DEVICE_PASSWORD" default:"pas.X!256$"`
}

func NewDeviceService(cfgReader *servercommon.ConfigurationReader,
	deviceRepository devicerepository.DeviceRepositorier,
	channelRepository channelrepository.ChannelRepositorier,
	commandRepository commandrepository.CommandRepositorier) (*DeviceService, error) {
	service := DeviceService{
		deviceRepository:  deviceRepository,
		channelRepository: channelRepository,
		commandRepository: commandRepository,
	}
	err := cfgReader.ReadConfiguration(&service, DEVICE_SERVICE_CONFIG_PATH)
	return &service, err
}

func (service *DeviceService) RegisterHandlers(r *http.ServeMux) {
	r.HandleFunc("POST /state", service.ESP8266Auth(service.ProcessPostState))
}
