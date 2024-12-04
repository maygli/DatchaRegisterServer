package devicerepository

import (
	"datcha/datamodel"

	"github.com/jinzhu/copier"
)

func getNewRepDevice(device *datamodel.Device) (*RepDevice, error) {
	repDevice := RepDevice{}
	err := copier.Copy(&repDevice, device)
	return &repDevice, err
}

func getNewDMDevice(repDevice *RepDevice) (*datamodel.Device, error) {
	dmDevice := datamodel.Device{}
	err := copier.Copy(&dmDevice, repDevice)
	return &dmDevice, err
}
