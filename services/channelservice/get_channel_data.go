package channelservice

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ChannelDataItem struct {
	Value     string
	Unit      string
	TimeStamp time.Time
}

func NewChannelDataItem(item *datamodel.ChannelData) ChannelDataItem {
	dataItem := ChannelDataItem{
		Value:     item.Data,
		Unit:      item.Unit,
		TimeStamp: item.CreatedAt,
	}
	return dataItem
}

func (service *ChannelService) getChannelDataHandle(w http.ResponseWriter, r *http.Request) {
	channelIdStr := r.PathValue(servercommon.CHANNEL_ID_KEY)
	if channelIdStr == "" {
		slog.Error("Error. channelId parameter notr present in path")
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	channelId, err := strconv.ParseUint(channelIdStr, 10, 64)
	if err != nil {
		slog.Error("Error. Can't parse channelId. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	chData, err := service.channelRepository.GetAllChannelData(uint(channelId))
	if err != nil {
		slog.Error("Error. Can't getChannleData. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	var res = make([]ChannelDataItem, len(chData), len(chData))
	for indx, dataItem := range chData {
		chDataItem := NewChannelDataItem(dataItem)
		res[indx] = chDataItem
	}
	resStr, err := json.Marshal(res)
	if err != nil {
		slog.Error("Error. Can't convert channel data to json. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	w.Write(resStr)
}
