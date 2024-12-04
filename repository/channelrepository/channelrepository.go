package channelrepository

import (
	"datcha/datamodel"
	"datcha/servercommon"

	"gorm.io/gorm"
)

type ChannelRepository struct {
	database *gorm.DB
	notifier servercommon.INotifier
}

func NewChannelRepository(db *gorm.DB, notifier servercommon.INotifier) (ChannelRepositorier, error) {
	rep := ChannelRepository{
		database: db,
		notifier: notifier,
	}
	err := rep.database.AutoMigrate(&RepChannel{})
	if err != nil {
		return &rep, err
	}
	err = rep.database.AutoMigrate(&RepChannelData{})
	return &rep, err
}

func (repo *ChannelRepository) FindChannel(channelId uint) (*datamodel.Channel, error) {
	repChannel := RepChannel{}
	result := repo.database.First(&repChannel, channelId)
	if result.Error != nil {
		return &datamodel.Channel{}, result.Error
	}
	return getNewDMChannel(&repChannel)
}

func (repo *ChannelRepository) FindChannelByName(deviceId uint, channelName string) (*datamodel.Channel, error) {
	channel := RepChannel{}
	result := repo.database.Where("deviceId = ? AND name=?", deviceId, channelName).First(&channel)
	if result.Error != nil {
		return &datamodel.Channel{}, result.Error
	}
	return getNewDMChannel(&channel)
}

func (repo *ChannelRepository) FindOrCreateChannelByName(deviceId uint, channelName string) (*datamodel.Channel, error) {
	channel := RepChannel{
		DeviceID: deviceId,
		Name:     channelName,
	}
	result := repo.database.FirstOrCreate(&channel, channel)
	if result.Error != nil {
		return &datamodel.Channel{}, result.Error
	}
	return getNewDMChannel(&channel)
}

func (repo *ChannelRepository) GetDeviceChannels(deviceId uint) ([]*datamodel.Channel, error) {
	channels := []RepChannel{}
	result := repo.database.Where("device_id=?", deviceId).Find(&channels)
	if result.Error != nil {
		return nil, result.Error
	}
	return getNewDMChannels(channels)
}
