package projectservice

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"fmt"
	"log/slog"
	"net/http"
)

func (server *ProjectService) newProjectsHandle(w http.ResponseWriter, r *http.Request) {
	projectData, err := server.getProjectData(r)
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
	slog.Debug(fmt.Sprintf("projectData=%v", projectData))
	project := toDmProject(projectData)
	user, ok := r.Context().Value(servercommon.USER_CONTEXT_KEY).(*datamodel.User)
	if !ok {
		slog.Error("context doesn't contains user data")
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	project.OwnerID = user.ID
	err = server.projectRepository.AddProject(&project)
	if err != nil {
		slog.Error("can't save project in the database. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	// save image
	projectFolder := servercommon.GetProjectFolder(project.ID)
	imageFileName := servercommon.PROJECT_IMAGE_FILENAME
	filePath, err := servercommon.SaveFormFile(r, "project_image", projectFolder, imageFileName)
	if err != nil {
		slog.Error("can't save project image. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	project.ImagePath = filePath
	err = server.projectRepository.UpdateProject(&project)
	if err != nil {
		slog.Error("can't update project. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
