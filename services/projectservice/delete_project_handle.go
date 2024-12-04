package projectservice

import (
	"datcha/servercommon"
	"log"
	"net/http"
	"os"
	"strconv"
)

func (service *ProjectService) deleteProjectsHandle(w http.ResponseWriter, r *http.Request) {
	projectIdStr := r.PathValue(servercommon.PROJECT_ID_KEY)
	if projectIdStr == "" {
		log.Printf("Error. projectId parameter notr present in path")
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	projectId, err := strconv.ParseUint(projectIdStr, 10, 64)
	if err != nil {
		log.Printf("Error. Can't parse projectId. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	err = service.projectRepository.DeleteProject(uint(projectId))
	if err != nil {
		log.Printf("Error. Can't delete project. Error: " + err.Error())
	}
	projectFolder := servercommon.GetProjectFolder(uint(projectId))
	err = os.RemoveAll(projectFolder)
	if err != nil {
		log.Printf("Error. Can't delete project folder. Error: " + err.Error())
	}
}
