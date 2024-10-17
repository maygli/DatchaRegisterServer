package commandserver

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"log"
	"net/http"
)

type CommandDto struct {
	Key       string `json:"key" schema:"key"`
	Value     string `json:"value" schema:"value"`
	DeviceId  uint   `json:"device_id" schema:"device_id"`
	ChannelId uint   `json:"channel_id" schema:"channel_id"`
}

func (server *CommandServer) putCommandHandle(w http.ResponseWriter, r *http.Request) {
	commandData := CommandDto{
		DeviceId:  servercommon.INVALID_ID,
		ChannelId: servercommon.INVALID_ID,
	}
	err := servercommon.ProcessBodyData(r, &commandData)
	if err != nil {
		log.Printf("Error. can't parse command data. Error=%s", err.Error())
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	if commandData.DeviceId == servercommon.INVALID_ID && commandData.ChannelId == servercommon.INVALID_ID {
		log.Printf("Error. Command data should contains device_id or channel_id field.")
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	cmd := datamodel.Command{
		Value:    commandData.Value,
		Key:      commandData.Key,
		DeviceID: commandData.DeviceId,
	}
	if cmd.DeviceID == servercommon.INVALID_ID {
		channel, err := server.repository.FindChannel(commandData.ChannelId)
		if err != nil {
			log.Printf("Error. Invalid channel id. Error=%s", err.Error())
			http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
			return
		}
		cmd.DeviceID = channel.DeviceID
	}
	err = server.repository.AddOrUpdateCommand(&cmd)
	if err != nil {
		log.Printf("Error. Can't write command. Error=%s", err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
