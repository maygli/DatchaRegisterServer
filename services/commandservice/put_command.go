package commandservice

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"log"
	"log/slog"
	"net/http"
)

type CommandDto struct {
	Key       string `json:"key" schema:"key"`
	Value     string `json:"value" schema:"value"`
	DeviceId  uint   `json:"device_id" schema:"device_id"`
	ChannelId uint   `json:"channel_id" schema:"channel_id"`
}

func (service *CommandService) putCommandHandle(w http.ResponseWriter, r *http.Request) {
	commandData := CommandDto{
		DeviceId:  servercommon.INVALID_ID,
		ChannelId: servercommon.INVALID_ID,
	}
	err := servercommon.ProcessBodyData(r, &commandData)
	if err != nil {
		slog.Error("Error. can't parse command data. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	if commandData.DeviceId == servercommon.INVALID_ID && commandData.ChannelId == servercommon.INVALID_ID {
		slog.Error("Error. Command data should contains device_id or channel_id field.")
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	cmd := datamodel.Command{
		Value:    commandData.Value,
		Key:      commandData.Key,
		DeviceID: commandData.DeviceId,
	}
	if cmd.DeviceID == servercommon.INVALID_ID {
		channel, err := service.channelRepository.FindChannel(commandData.ChannelId)
		if err != nil {
			slog.Error("Error. Invalid channel id. Error: " + err.Error())
			http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
			return
		}
		cmd.DeviceID = channel.DeviceID
	}
	err = service.cmdRepository.AddOrUpdateCommand(&cmd)
	if err != nil {
		log.Printf("Error. Can't write command. Error=%s", err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
