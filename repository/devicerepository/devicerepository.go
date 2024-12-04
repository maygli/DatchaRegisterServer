package devicerepository

import (
	"datcha/datamodel"
	"datcha/servercommon"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	database *gorm.DB
	notifier servercommon.INotifier
}

func NewDeviceRepository(db *gorm.DB, notifier servercommon.INotifier) (DeviceRepositorier, error) {
	rep := DeviceRepository{
		database: db,
		notifier: notifier,
	}
	err := rep.database.AutoMigrate(&RepDevice{})
	return &rep, err
}

func (repo *DeviceRepository) AddDevice(device *datamodel.Device) error {
	repoDevice, err := getNewRepDevice(device)
	if err != nil {
		return err
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(repoDevice)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *DeviceRepository) FindDevice(deviceId uint) (*datamodel.Device, error) {
	device := RepDevice{}
	result := repo.database.First(&device, deviceId)
	if result.Error != nil {
		return &datamodel.Device{}, result.Error
	}
	return getNewDMDevice(&device)
}
