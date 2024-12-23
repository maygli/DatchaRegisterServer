package profileservice

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"log/slog"
	"net/http"
)

func (service *ProfileService) putProfileHandle(w http.ResponseWriter, r *http.Request) {
	profileData := ProfileData{}
	err := servercommon.ProcessBodyData(r, &profileData)
	if err != nil {
		slog.Error("can't process profile data. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	if !profileData.isValid() {
		slog.Error("invalid profile data")
		http.Error(w, servercommon.ERROR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	user, ok := r.Context().Value(servercommon.USER_CONTEXT_KEY).(*datamodel.User)
	if !ok {
		slog.Error("context doesn't contains user data")
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	profile, err := service.profileRepository.FindProfileByUserId(user.ID)
	if err != nil {
		slog.Error("error get profile for user" + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	// save image
	projectFolder := servercommon.GetProfileFolder(profile.ID)
	imageFileName := servercommon.PROFILE_AVATAR_FILENAME
	filePath, err := servercommon.ReplaceFormFile(r, "avatar", projectFolder, imageFileName, profile.Avatar)
	if err != nil {
		slog.Error("can't save project image. Error: " + err.Error())
	} else {
		profileData.Avatar = filePath
	}
	mergeProfileData(profile, profileData)
	err = service.profileRepository.UpdateProfile(profile)
	if err != nil {
		slog.Error("can't update profile. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
}
