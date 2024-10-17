package repository

import (
	"datcha/datamodel"
	"datcha/servercommon"
)

func (repo *DatchaSqlRepository) AddOrUpdateCommand(command *datamodel.Command) error {
	var newCommand datamodel.Command
	result := repo.database.Where("device_id = ? AND key = ?", command.DeviceID, command.Key).First(&newCommand)
	if result.Error != nil {
		result = repo.database.Create(command)
	} else {
		newCommand.Value = command.Value
		command.ID = newCommand.ID
		repo.database.Save(&newCommand)
	}
	return result.Error
}

func (repo *DatchaSqlRepository) RemoveCommand(commandId uint) error {
	repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Delete(&datamodel.Command{}, commandId)
	return nil
}

func (repo *DatchaSqlRepository) GetDeviceCommands(deviceId uint) ([]datamodel.Command, error) {
	var commands []datamodel.Command
	result := repo.database.Where("device_id = ?", deviceId).Find(&commands)
	return commands, result.Error
}
