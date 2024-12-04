package projectrepository

import (
	"datcha/datamodel"

	"github.com/jinzhu/copier"
)

func getNewRepProject(project *datamodel.Project) (*RepProject, error) {
	repproject := RepProject{}
	err := copier.Copy(&repproject, project)
	return &repproject, err
}

func getNewDMProject(project *RepProject) (*datamodel.Project, error) {
	dmProject := datamodel.Project{}
	err := updateDMProject(project, &dmProject)
	return &dmProject, err
}

func updateDMProject(project *RepProject, dmProject *datamodel.Project) error {
	err := copier.Copy(&dmProject, project)
	return err
}

func getNewDMProjectsList(projects []RepProject) ([]*datamodel.Project, error) {
	result := make([]*datamodel.Project, 0, len(projects))
	for _, project := range projects {
		dmProject, err := getNewDMProject(&project)
		if err != nil {
			return result, err
		}
		result = append(result, dmProject)
	}
	return result, nil
}
