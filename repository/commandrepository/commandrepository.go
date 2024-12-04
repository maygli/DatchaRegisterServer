package commandrepository

import (
	"datcha/datamodel"
	"datcha/servercommon"

	"gorm.io/gorm"
)

type CommandRepository struct {
	database *gorm.DB
	notifier servercommon.INotifier
}

func NewCommandRepository(db *gorm.DB, notifier servercommon.INotifier) (CommandRepositorier, error) {
	rep := CommandRepository{
		database: db,
		notifier: notifier,
	}
	err := rep.database.AutoMigrate(&RepCommand{})
	if err != nil {
		return &rep, err
	}
	return &rep, nil
}

func (repo *CommandRepository) AddOrUpdateCommand(command *datamodel.Command) error {
	repCommand, err := getNewRepCommand(command)
	if err != nil {
		return err
	}
	result := repo.database.Where("device_id = ? AND key = ?", command.DeviceID, command.Key).First(&repCommand)
	if result.Error != nil {
		result = repo.database.Create(repCommand)
	} else {
		// TODO should convert from command but set saved id
		repCommand.Value = command.Value
		result = repo.database.Save(&repCommand)
	}
	return result.Error
}

func (repo *CommandRepository) RemoveCommand(commandId uint) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Delete(&datamodel.Command{}, commandId)
	return result.Error
}

func (repo *CommandRepository) GetDeviceCommands(deviceId uint) ([]*datamodel.Command, error) {
	var commands []RepCommand
	result := repo.database.Where("device_id = ?", deviceId).Find(&commands)
	if result.Error != nil {
		return []*datamodel.Command{}, result.Error
	}
	return getNewDMCommands(commands)
}
