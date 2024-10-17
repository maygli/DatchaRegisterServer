package main

import (
	"crypto/tls"
	"datcha/authentification"
	"datcha/channelserver"
	"datcha/commandserver"
	"datcha/deviceserver"
	"datcha/notifyserver"
	"datcha/pageserver"
	"datcha/projectserver"
	"datcha/repository"
	"datcha/servercommon"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"

	_ "github.com/lib/pq"
)

const (
	DATCHA_SERVER_CFG_FILE_NAME = "datcha.json"
	SERVER_HTTPS_MODE           = "https"
)

type DatchaServerConfiguration struct {
	ServerMode  string `json:"mode" env:"DATCHA_SERVER_MODE, overwrite, default=http"`
	LogFileName string `json:"log_file_name" env:"DATCHA_LOG_FILENAME, overwrite, default=datcha.log"`
}

func NewDatchaServerConfiguration() *DatchaServerConfiguration {
	config := DatchaServerConfiguration{}
	servercommon.ConfigureServer(&config, DATCHA_SERVER_CFG_FILE_NAME)
	return &config
}

func main() {
	serverCfg := NewDatchaServerConfiguration()
	logFile, err := os.OpenFile(serverCfg.LogFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	}

	mainServerRouter := mux.NewRouter()
	deviceServerRouter := mux.NewRouter()

	broker := notifyserver.NewNotifyServer()

	dbRepo, err := repository.NewSqlRepository(broker)
	if err != nil {
		log.Fatal("can't init repository. Error: ", err)
	}

	authServer, err := authentification.NewAuthServer(dbRepo)
	if authServer == nil {
		log.Fatal("can't init authentification server. Error: ", err)
	}

	deviceServer, err := deviceserver.NewDeviceServer(dbRepo)
	if err != nil {
		log.Fatal("can't init device server. Error: ", err)
	}

	channelServer := channelserver.NewChannelServer(authServer, dbRepo)
	commandServer := commandserver.NewCommandServer(authServer, dbRepo)

	token, err := deviceServer.GenerateToken(52)
	log.Printf("size=%d token=%s", len(token), token)

	//Authentification handler should registered before other handles
	authServer.RegisterHandlers(mainServerRouter)
	projectServer := projectserver.NewProjectServer(authServer, dbRepo)
	projectServer.RegisterHandlers(mainServerRouter)
	channelServer.RegisterHandlers(mainServerRouter)
	commandServer.RegisterHandlers(mainServerRouter)
	deviceServer.RegisterHandlers(deviceServerRouter)
	mainServerRouter.Handle("/notify", authServer.JwtAuth(http.HandlerFunc(broker.ServeHTTP))).Methods("GET")
	//It should be last one
	pageserver.RegisterHandlers(mainServerRouter)
	log.Println("Server started")
	go http.ListenAndServe(":6080", deviceServerRouter)
	if serverCfg.ServerMode == SERVER_HTTPS_MODE {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("crazydatcha.ru", "www.crazydatcha.ru"),
			Cache:      autocert.DirCache("./certs"), //Folder for storing certificates
		}
		server := &http.Server{
			Addr:    ":https",
			Handler: mainServerRouter,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		log.Fatal(server.ListenAndServeTLS("", "")) //Key and cert are coming from Let's Encrypt
	} else {
		http.ListenAndServe(":7080", mainServerRouter)
	}
}
