package projectserver

import (
	"datcha/servercommon"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func (server *ProjectServer) deleteProjectsHandle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	projectIdStr, ok := params[servercommon.PROJECT_ID_KEY]
	if !ok {
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
	err = server.repository.DeleteProject(uint(projectId))
	if err != nil {
		log.Printf("Error. Can't delete project. Error: " + err.Error())
	}
	projectFolder := servercommon.GetProjectFolder(uint(projectId))
	err = os.RemoveAll(projectFolder)
	if err != nil {
		log.Printf("Error. Can't delete project folder. Error: " + err.Error())
	}
}
