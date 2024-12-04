package channelrepository

import (
	"datcha/datamodel"
)

type ChannelRepositorier interface {
	FindChannel(channelId uint) (*datamodel.Channel, error)
	FindChannelByName(deviceId uint, channelName string) (*datamodel.Channel, error)
	FindOrCreateChannelByName(deviceId uint, channelName string) (*datamodel.Channel, error)
	GetDeviceChannels(deviceId uint) ([]*datamodel.Channel, error)

	AddChannelData(channelId uint, data string, unit string) error
	GetAllChannelData(channelId uint) ([]*datamodel.ChannelData, error)
	GetLastChannelData(channelId uint) (*datamodel.ChannelData, error)
}
