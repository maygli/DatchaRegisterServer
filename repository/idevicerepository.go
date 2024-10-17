package repository

import "datcha/datamodel"

type IDeviceRepository interface {
	AddDevice(project *datamodel.Device) error
	FindDevice(deviceId uint) (datamodel.Device, error)
	FindDeviceByName(projectId uint, name string) (datamodel.Device, error)
	RemoveDevice(id uint) error
}
