package projectservice

import (
	"datcha/servercommon"
	"fmt"
	"log/slog"
	"net/http"
)

func (service *ProjectService) putProjectsHandle(w http.ResponseWriter, r *http.Request) {
	projectData, err := service.getProjectData(r)
	if err != nil {
		slog.Error("can't get project data. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	err = projectData.IsValid()
	if err != nil {
		slog.Error("parsed project data is invalid. Error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	project, err := service.projectRepository.FindProject(projectData.ID)
	if err != nil {
		slog.Error(fmt.Sprintf("can't found project %d. Error: %s", project.ID, err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// save image
	projectFolder := servercommon.GetProjectFolder(projectData.ID)
	imageFileName := servercommon.PROJECT_IMAGE_FILENAME
	projectData.Image, err = servercommon.ReplaceFormFile(r, "project_image", projectFolder, imageFileName, project.ImagePath)
	if err != nil {
		slog.Error("can't save project image. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	slog.Debug(fmt.Sprintf("projectData=%v", projectData))
	mergeProjectData(project, projectData)
	err = service.projectRepository.UpdateProject(project)
	if err != nil {
		slog.Error("can't save project in the database. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
