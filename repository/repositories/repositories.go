package repositories

import (
	"datcha/repository/authrepository"
	"datcha/repository/channelrepository"
	"datcha/repository/commandrepository"
	"datcha/repository/devicerepository"
	"datcha/repository/projectcardrepository"
	"datcha/repository/projectrepository"
	"datcha/repository/repositorycommon"
	"datcha/servercommon"
)

type Repositories struct {
	AuthRep        authrepository.AuthRepositorier
	ProjectRep     projectrepository.ProjectRepositorier
	DeviceRep      devicerepository.DeviceRepositorier
	ChannelRep     channelrepository.ChannelRepositorier
	CommandRep     commandrepository.CommandRepositorier
	ProjectCardRep projectcardrepository.ProjectCardRepositorier
}

func NewRepositories(cfgReader *servercommon.ConfigurationReader, notifier servercommon.INotifier) (*Repositories, error) {
	reps := Repositories{}
	db, err := repositorycommon.NewDatabase(cfgReader)
	if err != nil {
		return nil, err
	}
	reps.AuthRep, err = authrepository.NewAuthRepository(db, notifier)
	if err != nil {
		return nil, err
	}
	reps.ProjectRep, err = projectrepository.NewProjectRepository(db, notifier)
	if err != nil {
		return nil, err
	}
	reps.DeviceRep, err = devicerepository.NewDeviceRepository(db, notifier)
	if err != nil {
		return nil, err
	}
	reps.ChannelRep, err = channelrepository.NewChannelRepository(db, notifier)
	if err != nil {
		return nil, err
	}
	reps.CommandRep, err = commandrepository.NewCommandRepository(db, notifier)
	if err != nil {
		return nil, err
	}
	reps.ProjectCardRep, err = projectcardrepository.NewProjectCardRepository(db, notifier)
	if err != nil {
		return nil, err
	}
	return &reps, nil
}
