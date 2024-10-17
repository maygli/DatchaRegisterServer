package projectserver

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"encoding/json"
	"log"
	"net/http"
)

func (server *ProjectServer) getProjectsHandle(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(servercommon.USER_CONTEXT_KEY).(datamodel.User)
	if !ok {
		log.Println("context doesn't contains user data")
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	log.Println("User name=", user)
	projects, err := server.repository.GetOwnerProjects(user.ID)
	if err != nil {
		log.Printf("Error. Can't get projects by user id(%d) . Error: %s", user.ID, err.Error())
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	var res = make([]ProjectData, 0, len(projects))
	for _, project := range projects {
		projectData := server.FromProject(project)
		res = append(res, projectData)
	}
	resStr, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error. Can't convert projects to json. Error: %s", err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	w.Write(resStr)
}
