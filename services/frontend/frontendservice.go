package frontend

import (
	"datcha/servercommon"
	"net/http"
)

const (
	FRONTEND_SERVER_CONFIGURATION_PATH = "$.frontendserver"
)

type FrontendService struct {
	FrontendFolder       string `json:"fronend_folder" env:"${SERVER_NAME}_FRONTEND_FOLDER" default:"data/static"`
	NotFoundRedirectPage string `json:"notfound_redirect_page" env:"${SERVER_NAME}_NOT_FOUND_REDIRECT_PAGE"`
	NotFoundFile         string `json:"notfound_file" env:"${SERVER_NAME}_NOT_FOUND_FILE" default:"index.html"`
	UsersDataFolder      string `json:"users_data_folder" env:"${SERVER_NAME}_USERS_DATA_FOLEDR" default:"users_data"`
}

func NewFrontEndService(cfgReader *servercommon.ConfigurationReader) (*FrontendService, error) {
	server := FrontendService{}
	err := cfgReader.ReadConfiguration(&server, FRONTEND_SERVER_CONFIGURATION_PATH)
	return &server, err
}

func testHandle(w http.ResponseWriter, r *http.Request) {
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
}

func (service FrontendService) RegisterHandlers(mux *http.ServeMux) {
	usersDataHandle := http.StripPrefix("/"+service.UsersDataFolder, http.FileServer(http.Dir(service.UsersDataFolder)))
	mux.Handle("GET /"+service.UsersDataFolder+"/", usersDataHandle)
	handler := http.StripPrefix("/", http.FileServer(http.Dir(service.FrontendFolder)))
	if service.NotFoundRedirectPage != "" || service.NotFoundFile != "" {
		handler = service.fsMiddleware(handler)
	}
	mux.Handle("GET /", handler)
}
