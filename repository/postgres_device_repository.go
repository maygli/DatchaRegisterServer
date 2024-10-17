package repository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"
)

func (repo *DatchaSqlRepository) AddDevice(device *datamodel.Device) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(device)
	if result.Error != nil {
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	return nil
}

func (repo *DatchaSqlRepository) FindDevice(deviceId uint) (datamodel.Device, error) {
	device := datamodel.Device{}
	result := repo.database.First(&device, deviceId)
	return device, result.Error
}

func (repo *DatchaSqlRepository) FindDeviceByName(projectId uint, name string) (datamodel.Device, error) {
	return datamodel.Device{}, nil
}

func (repo *DatchaSqlRepository) RemoveDevice(id uint) error {
	repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Delete(&datamodel.Device{}, id)
	return nil
}
