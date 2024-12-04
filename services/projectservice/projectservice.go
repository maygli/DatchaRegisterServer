package projectservice

import (
	"datcha/datamodel"
	"datcha/repository/channelrepository"
	"datcha/repository/devicerepository"
	"datcha/repository/projectcardrepository"
	"datcha/repository/projectrepository"
	"datcha/servercommon"
	"datcha/services/authservice/authinternalservice"
	"errors"
	"net/http"
	"time"
)

const (
	MIN_PROJECT_NAME_LEN int = 4
)

type ProjectService struct {
	projectRepository     projectrepository.ProjectRepositorier
	channelRepository     channelrepository.ChannelRepositorier
	deviceRepository      devicerepository.DeviceRepositorier
	projectCardRepository projectcardrepository.ProjectCardRepositorier
	authService           *authinternalservice.AuthInternalService
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

func NewProjectService(rep projectrepository.ProjectRepositorier,
	devRep devicerepository.DeviceRepositorier,
	chanRep channelrepository.ChannelRepositorier,
	cardRep projectcardrepository.ProjectCardRepositorier,
	authservice *authinternalservice.AuthInternalService) *ProjectService {
	service := ProjectService{
		projectRepository:     rep,
		deviceRepository:      devRep,
		channelRepository:     chanRep,
		projectCardRepository: cardRep,
		authService:           authservice,
	}
	return &service
}

func (service *ProjectService) getProjectData(r *http.Request) (ProjectData, error) {
	projectData := ProjectData{}
	err := servercommon.ProcessBodyData(r, &projectData)
	return projectData, err
}

func (service *ProjectService) ToProject(projectData ProjectData) datamodel.Project {
	project := datamodel.Project{
		Name: projectData.Name,
	}
	return project
}

func (service *ProjectService) fromCardItem(cardItem *datamodel.ProjectCardItem) (ChannelLastValue, error) {
	val := ChannelLastValue{
		ChannelId: cardItem.ChannelID,
	}
	channel, err := service.channelRepository.FindChannel(cardItem.ChannelID)
	if err != nil {
		return val, err
	}
	val.ChannelName = channel.Name
	lastData, err := service.channelRepository.GetLastChannelData(cardItem.ChannelID)
	if err == nil {
		val.LatsValue = lastData.Data
		val.Unit = lastData.Unit
		val.CreatedAt = lastData.CreatedAt
	}
	device, err := service.deviceRepository.FindDevice(channel.DeviceID)
	if err == nil {
		val.DeviceId = channel.DeviceID
		val.DeviceName = device.Name
	}
	return val, nil
}

func (service *ProjectService) FromProject(project *datamodel.Project) ProjectData {
	projectData := ProjectData{
		ID:      project.ID,
		Name:    project.Name,
		Image:   project.ImagePath,
		OwnerId: project.OwnerID,
	}
	projectData.ProjectCardsItems = make([]ChannelLastValue, 0)
	cardItems, err := service.projectCardRepository.GetProjectCardItems(project.ID)
	if err == nil {
		for _, cardItem := range cardItems {
			lastVal, err := service.fromCardItem(cardItem)
			if err == nil {
				projectData.ProjectCardsItems = append(projectData.ProjectCardsItems, lastVal)
			}
		}
	}
	return projectData
}

func (service *ProjectService) RegisterHandlers(r *http.ServeMux) {
	r.Handle("GET /projects", service.authService.JwtAuth(http.HandlerFunc(service.getProjectsHandle)))
	r.Handle("POST /projects", service.authService.JwtAuth(http.HandlerFunc(service.newProjectsHandle)))
	r.Handle("DELETE /project/{"+servercommon.PROJECT_ID_KEY+"}", service.authService.JwtAuth(http.HandlerFunc(service.deleteProjectsHandle)))
}
