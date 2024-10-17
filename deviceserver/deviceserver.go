package deviceserver

import (
	"datcha/repository"
	"datcha/servercommon"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

const (
	DEVICE_SERVER_CONFIG_FILE_NAME string = "deviceserverconfig.json"
)

type DeviceServer struct {
	MapMutex       sync.Mutex
	repository     repository.IRepository
	Issuer         string `json:"issuer" env:"DATCHA_DEVICE_ISSUER, overwrite, default=crazydatchahome"`
	StateSubject   string `json:"state_subject" env:"DATCHA_DEVICE_STATE_SUBJECT, overwrite, default=devicesstate"`
	StateSecretKey string `json:"state_secret" env:"DATCHA_DEVICE_STATE_SECRET, overwrite, default=,g1XPz){ubdSK8{_(%rra[&t[rWSpEVV"`
	DeviceUser     string `json:"device_user" env:"DATCHA_DEVICE_USER, overwrite, default=device"`
	DevicePassword string `json:"device_password" env:"DATCHA_DEVICE_PASSWORD, overwrite, default=pas.X!256$"`
}

func NewDeviceServer(rep repository.IRepository) (*DeviceServer, error) {
	server := DeviceServer{}
	server.repository = rep
	err := servercommon.ConfigureServer(&server, DEVICE_SERVER_CONFIG_FILE_NAME)
	return &server, err
}

func (server *DeviceServer) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/state", server.ESP8266Auth(server.ProcessPostState)).Methods(http.MethodPost)
}
