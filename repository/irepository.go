package repository

import (
	"datcha/datamodel"
)

type IRepository interface {
	AddUser(user *datamodel.User) error
	FindUser(id uint) (datamodel.User, error)
	FindUserByName(name string) (datamodel.User, error)
	FindUserByEmail(email string) (datamodel.User, error)
	FindUserByNameEmail(name string) (datamodel.User, error)
	DeleteUser(id uint) error
	UpdateAccountStatus(id uint, status datamodel.AccountStatus) error

	AddProject(project *datamodel.Project) error
	RemoveProject(projectId uint) error
	FindProject(projectId uint) (datamodel.Project, error)
	UpdateProject(project *datamodel.Project) error
	DeleteProject(projectId uint) error
	GetOwnerProjects(userId uint) ([]datamodel.Project, error)

	AddDevice(device *datamodel.Device) error
	FindDevice(deviceId uint) (datamodel.Device, error)

	FindChannel(channelId uint) (datamodel.Channel, error)
	FindChannelByName(deviceId uint, channelName string) (datamodel.Channel, error)
	FindOrCreateChannelByName(deviceId uint, channelName string) (datamodel.Channel, error)
	GetDeviceChannels(deviceId uint) ([]datamodel.Channel, error)

	AddChannelData(channelId uint, data string, unit string) error
	GetAllChannelData(channelId uint) ([]datamodel.ChannelData, error)
	GetLastChannelData(channelId uint) (datamodel.ChannelData, error)

	AddProjectCardItem(cardItem *datamodel.ProjectCardItem) error
	RemoveProjectCardItem(projectCardItemId uint) error
	GetProjectCardItems(projectId uint) ([]datamodel.ProjectCardItem, error)

	AddOrUpdateCommand(command *datamodel.Command) error
	RemoveCommand(commandId uint) error
	GetDeviceCommands(deviceId uint) ([]datamodel.Command, error)
}
