package deviceservice

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"

	//	"io/ioutil"
	"net/http"
)

const STATE_CHANNEL string = "state"
const UNIT_SUFFIX string = "_unit"

func convertToString(data interface{}) (string, error) {
	strData, ok := data.(string)
	if !ok {
		floatData, ok := data.(float64)
		if ok {
			strData = strconv.FormatFloat(floatData, 'g', 15, 64)
		} else {
			boolData, ok := data.(bool)
			if ok {
				strData = strconv.FormatBool(boolData)
			} else {
				return "", errors.New("unsupported type")
			}
		}
	}
	return strData, nil
}

func (service *DeviceService) getCommandsMap(deviceId uint) (map[string]string, error) {
	cmds, err := service.commandRepository.GetDeviceCommands(deviceId)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, cmd := range cmds {
		result[cmd.Key] = cmd.Value
	}
	return result, nil
}

func (service *DeviceService) ProcessPostState(w http.ResponseWriter, r *http.Request) {
	var deviceInfo map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&deviceInfo)
	if err != nil {
		slog.Error("Error. Can't read request body: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	device, ok := r.Context().Value(servercommon.DEVICE_CONTEXT_KEY).(*datamodel.Device)
	if !ok {
		slog.Error("context doesn't contains device data")
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	for channelId, value := range deviceInfo {
		if strings.HasSuffix(channelId, UNIT_SUFFIX) {
			continue
		}
		unitKey := channelId + UNIT_SUFFIX
		unit, _ := deviceInfo[unitKey]
		channel, err := service.channelRepository.FindOrCreateChannelByName(device.ID, channelId)
		if err != nil {
			errMsg := fmt.Sprintf("Error. Can't find or create channel deviceId=%d, channelId=%s. Error '%s'", device.ID, channelId, err.Error())
			slog.Error(errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		valueStr, err := convertToString(value)
		if err != nil {
			slog.Error("Can't convert received value to string. Error=" + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		unitStr := ""
		if unit != nil {
			unitStr, ok = unit.(string)
			if !ok {
				msg := "Can't convert unit value to string"
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}
		err = service.channelRepository.AddChannelData(channel.ID, valueStr, unitStr)
		if err != nil {
			msg := "Can't write channel data. Error:" + err.Error()
			slog.Error(msg)
			http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
			return
		}
		fmt.Println(channel)
	}
	respMap, err := service.getCommandsMap(device.ID)
	if err != nil {
		msg := "Can't generates respons map. Error:" + err.Error()
		log.Println(msg)
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	resStr, err := json.Marshal(respMap)
	if err != nil {
		msg := "Can't generates respons json. Error:" + err.Error()
		log.Println(msg)
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	w.Write(resStr)
	//	fmt.Fprintf(w, resp)
}
