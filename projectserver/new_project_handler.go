package projectserver

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (server *ProjectServer) newProjectsHandle(w http.ResponseWriter, r *http.Request) {
	projectData, err := server.getProjectData(r)
	if err != nil {
		log.Println("can't get project data. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	err = projectData.IsValid()
	if err != nil {
		log.Println("parsed project data is invalid. Error: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("projectData=", projectData)
	project := server.ToProject(projectData)
	user, ok := r.Context().Value(servercommon.USER_CONTEXT_KEY).(datamodel.User)
	if !ok {
		log.Println("context doesn't contains user data")
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	project.Owner = &user
	err = server.repository.AddProject(&project)
	if err != nil {
		log.Println("can't save project in the database. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	// in your case file would be fileupload
	file, header, err := r.FormFile("project_image")
	if err == nil {
		defer file.Close()
		projectFolder := servercommon.GetProjectFolder(project.ID)
		err = servercommon.CreateFolder(projectFolder)
		if err != nil {
			// Do nothing in case it is impossible to create project
			log.Println("can't create project folder ", projectFolder)
			return
		}
		imageFileName := filepath.Join(projectFolder, servercommon.PROJECT_IMAGE_FILENAME+filepath.Ext(header.Filename))
		log.Printf("File name %s\n", imageFileName)
		// Copy the file data to my buffer
		f, err := os.OpenFile(imageFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("can't open file to save iamge. FileName=", imageFileName)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			log.Println("can't save iamge. FileName=", imageFileName)
			return
		}
		project.ImagePath = imageFileName
		err = server.repository.UpdateProject(&project)
		if err != nil {
			log.Println("can't update project. Error: ", err.Error())
			http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
			return
		}
	}
}
