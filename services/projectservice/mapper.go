package projectservice

import "datcha/datamodel"

func toDmProject(projectData ProjectData) datamodel.Project {
	project := datamodel.Project{
		ID:        projectData.ID,
		OwnerID:   projectData.OwnerId,
		Name:      projectData.Name,
		ImagePath: projectData.Image,
	}
	return project
}

func mergeProjectData(project *datamodel.Project, data ProjectData) {
	if data.Name != "" {
		project.Name = data.Name
	}
	if data.Image != "" && data.Image != "null" {
		project.ImagePath = data.Image
	}
}
