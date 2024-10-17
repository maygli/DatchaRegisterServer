package repository

import "datcha/datamodel"

type IChannelRepository interface {
	AddChannel(deviceId uint, channel *datamodel.Channel) error
	FindChannel(deviceId uint, name string) (datamodel.Channel, error)
	FindOrCreateChannel(deviceId uint, name string) (datamodel.Channel, error)
	RemoveChannel(id uint) error
}
