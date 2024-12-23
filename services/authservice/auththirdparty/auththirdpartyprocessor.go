package auththirdparty

import (
	"datcha/datamodel"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	URL_REDIRECT_ERROR   = "/"
	URL_REDIRECT_SUCCESS = "/third_party"
)

func (service *AuthThirdPartyService) convertUser(thirdPartyUser AuthThirdPartyUser) *datamodel.User {
	user := datamodel.User{
		Name:      thirdPartyUser.GetName(),
		Email:     thirdPartyUser.Email,
		AccStatus: datamodel.AS_CONFIRMED,
	}
	user.SetServiceUserId(thirdPartyUser.Service, thirdPartyUser.UserId)
	return &user
}

func (service *AuthThirdPartyService) findOrAddUser(thirdPartyUser AuthThirdPartyUser) *datamodel.User {
	var dmUser *datamodel.User
	var err error
	if thirdPartyUser.Email != "" {
		dmUser, err = service.authRepository.FindUserByEmail(thirdPartyUser.Email)
		if err != nil {
			dmUser = nil
		}
	}
	if dmUser == nil {
		dmUser, err = service.authRepository.FindUserByServiceId(thirdPartyUser.UserId, thirdPartyUser.Service)
		if err != nil {
			dmUser = nil
		}
	}
	if dmUser == nil {
		// TODO if user with name already exists should generate unique name
		dmUser = service.convertUser(thirdPartyUser)
		err = service.authRepository.AddUser(dmUser)
		if err != nil {
			slog.Error("Can't create new user. Error=" + err.Error())
		}
	}
	return dmUser
}

func (service *AuthThirdPartyService) updateUserByThirdPartyData(dmUser *datamodel.User, thirdPartyUser AuthThirdPartyUser) {
	if !dmUser.HasServiceId(thirdPartyUser.Service) {
		dmUser.SetServiceUserId(thirdPartyUser.Service, thirdPartyUser.UserId)
		err := service.authRepository.UpdateServiceUserId(dmUser.ID, thirdPartyUser.Service, thirdPartyUser.UserId)
		if err != nil {
			slog.Error("third party auth. can't update user service id. Error=" + err.Error())
		}
	}
	if dmUser.Email == "" && thirdPartyUser.Email != "" {
		dmUser.Email = thirdPartyUser.Email
		err := service.authRepository.UpdateEmail(dmUser.ID, thirdPartyUser.Email)
		if err != nil {
			slog.Error("third party auth. can't update user email. Error=" + err.Error())
		}
	}
}

func (service *AuthThirdPartyService) updateProfile(dmProfile *datamodel.UserProfile, thirdPartyUser AuthThirdPartyUser) error {
	if dmProfile.Name == "" && thirdPartyUser.GetName() != "" {
		err := service.profileRepository.SetName(dmProfile.ID, thirdPartyUser.GetName())
		if err != nil {
			slog.Error("third party auth. can't update profile user name. Error=" + err.Error())
		}
	}
	if dmProfile.Avatar == "" && thirdPartyUser.Avatar != "" {
		err := service.profileRepository.SetAvatar(dmProfile.ID, thirdPartyUser.Avatar)
		if err != nil {
			slog.Error("third party auth. can't update profile avatar. Error=" + err.Error())
		}
	}
	return nil
}

func (service *AuthThirdPartyService) updateProfileByThirdPartyData(userId uint, thirdPartyUser AuthThirdPartyUser) error {
	profile, err := service.profileRepository.FindProfileByUserId(userId)
	if err != nil || profile == nil {
		profile := &datamodel.UserProfile{
			UserID: userId,
			Avatar: thirdPartyUser.Avatar,
			Name:   thirdPartyUser.GetName(),
		}
		return service.profileRepository.AddProfile(profile)
	}
	return service.updateProfile(profile, thirdPartyUser)
}

func (service *AuthThirdPartyService) ProcessSuceess(user any, w http.ResponseWriter, r *http.Request) error {
	thirdPartyUser, isOk := user.(AuthThirdPartyUser)
	if !isOk {
		slog.Error("can't process user data. Unsupported type")
		http.Redirect(w, r, URL_REDIRECT_ERROR, http.StatusSeeOther)
		return errors.New("unsupport user type")
	}
	dmUser := service.findOrAddUser(thirdPartyUser)
	if dmUser == nil {
		return errors.New("can't process third party user")
	}
	service.updateUserByThirdPartyData(dmUser, thirdPartyUser)
	service.updateProfileByThirdPartyData(dmUser.ID, thirdPartyUser)
	slog.Debug(fmt.Sprintf("Get or registered user by google. User=%+v", dmUser))
	service.intServive.SendAuthCoockie(w, dmUser.ID)
	err := service.intServive.WriteAutorizationHeader(w, dmUser.ID)
	if err != nil {
		slog.Error("can't write authoorization header. Error: " + err.Error())
		http.Redirect(w, r, URL_REDIRECT_ERROR, http.StatusSeeOther)
		return err
	}
	http.Redirect(w, r, URL_REDIRECT_SUCCESS, http.StatusSeeOther)
	return nil
}

func (service *AuthThirdPartyService) ProcessError(err error, w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, URL_REDIRECT_ERROR, http.StatusSeeOther)
	return nil
}
