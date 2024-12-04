package channelrepository

import (
	"datcha/datamodel"

	"github.com/jinzhu/copier"
)

func getNewDMChannel(repChannel *RepChannel) (*datamodel.Channel, error) {
	dmChannel := datamodel.Channel{}
	err := copier.Copy(&dmChannel, repChannel)
	return &dmChannel, err
}

func getNewDMChannels(channels []RepChannel) ([]*datamodel.Channel, error) {
	result := make([]*datamodel.Channel, len(channels))
	for indx, channel := range channels {
		dmChannel, err := getNewDMChannel(&channel)
		if err != nil {
			return result, err
		}
		result[indx] = dmChannel
	}
	return result, nil
}

func getNewDMChannelData(data *RepChannelData) (*datamodel.ChannelData, error) {
	dmData := datamodel.ChannelData{}
	err := copier.Copy(&dmData, data)
	return &dmData, err
}

func getNewDMChannelDatas(datas []RepChannelData) ([]*datamodel.ChannelData, error) {
	result := make([]*datamodel.ChannelData, len(datas))
	for indx, data := range datas {
		dmData, err := getNewDMChannelData(&data)
		if err != nil {
			return result, err
		}
		result[indx] = dmData
	}
	return result, nil
}
