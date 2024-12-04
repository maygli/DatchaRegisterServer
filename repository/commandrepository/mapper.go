package commandrepository

import (
	"datcha/datamodel"

	"github.com/jinzhu/copier"
)

func getNewRepCommand(command *datamodel.Command) (*RepCommand, error) {
	repCommand := RepCommand{}
	err := copier.Copy(&repCommand, command)
	return &repCommand, err
}

func getNewDMCommand(command *RepCommand) (*datamodel.Command, error) {
	dmCommand := datamodel.Command{}
	err := copier.Copy(&dmCommand, command)
	return &dmCommand, err
}

func getNewDMCommands(commands []RepCommand) ([]*datamodel.Command, error) {
	result := make([]*datamodel.Command, len(commands))
	for indx, command := range commands {
		dmCommand, err := getNewDMCommand(&command)
		if err != nil {
			return result, err
		}
		result[indx] = dmCommand
	}
	return result, nil
}
