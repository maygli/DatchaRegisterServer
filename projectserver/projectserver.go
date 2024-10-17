package projectserver

import (
	"datcha/authentification"
	"datcha/datamodel"
	"datcha/repository"
	"datcha/servercommon"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	MIN_PROJECT_NAME_LEN int = 4
)

type ProjectServer struct {
	repository repository.IRepository
	authServer *authentification.AuthServer
}

type ChannelLastValue struct {
	ChannelId   uint      `json:"channel_id" schema:"channel_id"`
	DeviceId    uint      `json:"device_id" schema:"device_id"`
	DeviceName  string    `json:"device_name" schema:"device_name"`
	ChannelName string    `json:"channel_name" schema:"channel_name"`
	LatsValue   string    `json:"last_value" schema:"last_value"`
	Unit        string    `json:"unit" schema:"unit"`
	CreatedAt   time.Time `json:"timestamp" schema:"timestamp"`
}

type ProjectData struct {
	ID                uint               `json:"project_id" schema:"project_id"`
	Name              string             `json:"project_name" schema:"project_name"`
	Image             string             `json:"project_image" schema:"project_image"`
	OwnerId           uint               `json:"owner_id" schema:"owner_id"`
	ProjectCardsItems []ChannelLastValue `json:"card_items" schema:"card_items"`
}

func (projectData ProjectData) IsValid() error {
	if len(projectData.Name) < MIN_PROJECT_NAME_LEN {
		return errors.New(servercommon.ERROR_PROJECT_NAME_TOO_SHORT)
	}
	return nil
}

func NewProjectServer(authServer *authentification.AuthServer, rep repository.IRepository) *ProjectServer {
	server := ProjectServer{}
	server.repository = rep
	server.authServer = authServer
	return &server
}

func (server *ProjectServer) getProjectData(r *http.Request) (ProjectData, error) {
	projectData := ProjectData{}
	err := servercommon.ProcessBodyData(r, &projectData)
	return projectData, err
}

func (server *ProjectServer) ToProject(projectData ProjectData) datamodel.Project {
	project := datamodel.Project{
		Name: projectData.Name,
	}
	return project
}

func (server *ProjectServer) fromCardItem(cardItem datamodel.ProjectCardItem) (ChannelLastValue, error) {
	val := ChannelLastValue{
		ChannelId: cardItem.ChannelID,
		//		DeviceId:
		//		DeviceName:
		//		ChannelName: cardItem.
		//		LatsValue   string    `json:"last_value" schema:"last_value"`
		//		Unit        string    `json:"unit" schema:"unit"`
		//		CreatedAt: cardItem.CreatedAt,
	}
	channel, err := server.repository.FindChannel(cardItem.ChannelID)
	if err != nil {
		return val, err
	}
	val.ChannelName = channel.Name
	lastData, err := server.repository.GetLastChannelData(cardItem.ChannelID)
	if err == nil {
		val.LatsValue = lastData.Data
		val.Unit = lastData.Unit
		val.CreatedAt = lastData.CreatedAt
	}
	device, err := server.repository.FindDevice(channel.DeviceID)
	if err == nil {
		val.DeviceId = channel.DeviceID
		val.DeviceName = device.Name
	}
	return val, nil
}

func (server *ProjectServer) FromProject(project datamodel.Project) ProjectData {
	projectData := ProjectData{
		ID:      project.ID,
		Name:    project.Name,
		Image:   project.ImagePath,
		OwnerId: project.OwnerID,
	}
	projectData.ProjectCardsItems = make([]ChannelLastValue, 0)
	for _, cardItem := range project.ProjectCardItems {
		lastVal, err := server.fromCardItem(cardItem)
		if err == nil {
			projectData.ProjectCardsItems = append(projectData.ProjectCardsItems, lastVal)
		}
	}
	return projectData
}

func (server *ProjectServer) RegisterHandlers(r *mux.Router) {
	r.Handle("/projects", server.authServer.JwtAuth(http.HandlerFunc(server.getProjectsHandle))).Methods("GET")
	r.Handle("/projects", server.authServer.JwtAuth(http.HandlerFunc(server.newProjectsHandle))).Methods("POST")
	r.Handle("/project/{"+servercommon.PROJECT_ID_KEY+"}", server.authServer.JwtAuth(http.HandlerFunc(server.deleteProjectsHandle))).Methods("DELETE")
}
