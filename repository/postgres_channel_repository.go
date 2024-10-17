package repository

import "datcha/datamodel"

func (repo *DatchaSqlRepository) FindChannel(channelId uint) (datamodel.Channel, error) {
	channel := datamodel.Channel{}
	result := repo.database.First(&channel, channelId)
	return channel, result.Error
}

func (repo *DatchaSqlRepository) FindChannelByName(deviceId uint, channelName string) (datamodel.Channel, error) {
	channel := datamodel.Channel{}
	result := repo.database.Where("deviceId = ? AND name=?", deviceId, channelName).First(&channel)
	return channel, result.Error
}

func (repo *DatchaSqlRepository) FindOrCreateChannelByName(deviceId uint, channelName string) (datamodel.Channel, error) {
	channel := datamodel.Channel{
		DeviceID: deviceId,
		Name:     channelName,
	}
	result := repo.database.FirstOrCreate(&channel, channel)
	return channel, result.Error
}

func (repo *DatchaSqlRepository) GetDeviceChannels(deviceId uint) ([]datamodel.Channel, error) {
	channels := []datamodel.Channel{}
	result := repo.database.Where("device_id=?", deviceId).Find(&channels)
	return channels, result.Error
}
