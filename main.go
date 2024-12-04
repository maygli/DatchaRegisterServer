package main

import (
	"crypto/tls"
	"datcha/repository/repositories"
	"datcha/servercommon"
	"datcha/serverlogger"
	"datcha/services/authservice"
	"datcha/services/channelservice"
	"datcha/services/commandservice"
	"datcha/services/deviceservice"
	"datcha/services/frontend"
	"datcha/services/notifyservice"
	"datcha/services/projectservice"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

const (
	SERVER_CFG_FILE_NAME = "data/config/server.json"
	SERVER_HTTPS_MODE    = "https"
	SERVER_NAME_KEY      = "SERVER_NAME"
	SERVER_NAME          = "TRAVEL"
)

type TravelServerConfiguration struct {
	ServerMode      string `json:"mode" env:"${SERVER_NAME}_SERVER_MODE" default:"http"`
	UsersDataFolder string `json:"users_data_folder" env:"${SERVER_NAME}_USERS_DATA_FOLDER" default:"data/users_data"`
	Port            int    `json:"port" env:"${SERVER_NAME}_CONFIG_PORT" default:"7080"`
}

func NewTravelServerConfiguration(cfgReader *servercommon.ConfigurationReader) (*TravelServerConfiguration, error) {
	config := TravelServerConfiguration{}
	err := cfgReader.ReadConfiguration(&config, "$")
	return &config, err
}

func main() {
	dict := make(map[string]string)
	dict[SERVER_NAME_KEY] = SERVER_NAME
	cfgReader, err := servercommon.NewConfigurationReader(SERVER_CFG_FILE_NAME, dict)
	if err != nil {
		slog.Error(fmt.Sprintf("can't get server configuration. Error: %s", err.Error()))
		os.Exit(1)
	}
	serverCfg, err := NewTravelServerConfiguration(cfgReader)
	if err != nil {
		slog.Error(fmt.Sprintf("can't get server configuration. Error: %s", err.Error()))
		os.Exit(1)
	}
	err = serverlogger.InitLogger(cfgReader)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	fmt.Println(serverCfg)
	mainServerRouter := http.NewServeMux()
	frontservice, err := frontend.NewFrontEndService(cfgReader)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	frontservice.RegisterHandlers(mainServerRouter)
	notifier := notifyservice.NewNotifyService()
	notifier.RegisterHandlers(mainServerRouter)
	reps, err := repositories.NewRepositories(cfgReader, notifier)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	authsrv, err := authservice.NewAuthService(cfgReader, reps.AuthRep)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	authsrv.RegisterHandlers(mainServerRouter)
	projsrv := projectservice.NewProjectService(reps.ProjectRep, reps.DeviceRep,
		reps.ChannelRep, reps.ProjectCardRep, authsrv.InternalService)
	projsrv.RegisterHandlers(mainServerRouter)
	deviceSrv, err := deviceservice.NewDeviceService(cfgReader, reps.DeviceRep, reps.ChannelRep, reps.CommandRep)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	deviceSrv.RegisterHandlers(mainServerRouter)
	channelSrv := channelservice.NewChannelService(authsrv.InternalService, reps.ChannelRep)
	channelSrv.RegisterHandlers(mainServerRouter)
	commandSrv := commandservice.NewCommandService(authsrv.InternalService, reps.CommandRep, reps.ChannelRep)
	commandSrv.RegisterHandlers(mainServerRouter)
	slog.Info("Server started")
	wrappedRouter := serverlogger.LoggerWrap(mainServerRouter)
	if serverCfg.ServerMode == SERVER_HTTPS_MODE {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("crazydatcha.ru", "www.crazydatcha.ru"),
			Cache:      autocert.DirCache("./certs"), //Folder for storing certificates
		}
		server := &http.Server{
			Addr:    ":https",
			Handler: wrappedRouter,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		slog.Error(server.ListenAndServeTLS("", "").Error()) //Key and cert are coming from Let's Encrypt
	} else {
		slog.Error(http.ListenAndServe(":80", wrappedRouter).Error())
	}
}
