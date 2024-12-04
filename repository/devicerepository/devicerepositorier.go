package devicerepository

import "datcha/datamodel"

type DeviceRepositorier interface {
	AddDevice(device *datamodel.Device) error
	FindDevice(deviceId uint) (*datamodel.Device, error)
}
