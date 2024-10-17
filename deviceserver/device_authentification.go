package deviceserver

import (
	"context"
	"datcha/servercommon"
	"fmt"
	"log"
	"net/http"
)

const (
	DEVICE_HEADER string = "Device"
)

func (server *DeviceServer) ESP8266Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request")
		userName, password, ok := r.BasicAuth()
		if !ok {
			// Workaround for esp8266. It's fail with strange error. Don't send
			// authenticae header - just return error.
			//		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			log.Println("Error. Can't get basic auth info from device request")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if userName != server.DeviceUser || password != server.DevicePassword {
			log.Println("Error can't get basic auth info from device request. userName=" + userName + " password=" + password)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		deviceHeader := r.Header.Get(DEVICE_HEADER)
		if deviceHeader == "" {
			log.Println("Error.Device request dosn't contains device header")
			http.Error(w, "no device header", http.StatusBadRequest)
			return
		}
		deviceId, err := server.ParseToken(deviceHeader)
		if err != nil {
			log.Println("Error. Can't parse device header. Err=", err.Error())
			http.Error(w, "can't parse device header", http.StatusBadRequest)
			return
		}
		device, err := server.repository.FindDevice(deviceId)
		if err != nil {
			log.Println("Error. Can't find device with id. Err=", err.Error())
			http.Error(w, "can't find device", http.StatusNotFound)
			return
		}
		newContext := context.WithValue(r.Context(), servercommon.DEVICE_CONTEXT_KEY, device)
		next.ServeHTTP(w, r.WithContext(newContext))
	})
}
