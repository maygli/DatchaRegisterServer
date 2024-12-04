package commandrepository

import "datcha/datamodel"

type CommandRepositorier interface {
	AddOrUpdateCommand(command *datamodel.Command) error
	RemoveCommand(commandId uint) error
	GetDeviceCommands(deviceId uint) ([]*datamodel.Command, error)
}
