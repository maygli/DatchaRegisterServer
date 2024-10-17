package channelserver

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ChannelDataItem struct {
	Value     string
	Unit      string
	TimeStamp time.Time
}

func NewChannelDataItem(item datamodel.ChannelData) ChannelDataItem {
	dataItem := ChannelDataItem{
		Value:     item.Data,
		Unit:      item.Unit,
		TimeStamp: item.CreatedAt,
	}
	return dataItem
}

func (server *ChannelServer) getChannelDataHandle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	channelIdStr, ok := params[servercommon.CHANNEL_ID_KEY]
	if !ok {
		log.Printf("Error. channelId parameter notr present in path")
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	channelId, err := strconv.ParseUint(channelIdStr, 10, 64)
	if err != nil {
		log.Printf("Error. Can't parse channelId. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	chData, err := server.repository.GetAllChannelData(uint(channelId))
	if err != nil {
		log.Printf("Error. Can't getChannleData. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	var res = make([]ChannelDataItem, 0, len(chData))
	for _, dataItem := range chData {
		chDataItem := NewChannelDataItem(dataItem)
		res = append(res, chDataItem)
	}
	resStr, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error. Can't channel data to json. Error: %s", err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	w.Write(resStr)
}
