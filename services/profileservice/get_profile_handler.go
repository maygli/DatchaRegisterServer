package profileservice

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (service *ProfileService) getProfileHandle(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(servercommon.USER_CONTEXT_KEY).(*datamodel.User)
	if !ok {
		slog.Error("context doesn't contains user data")
		http.Error(w, servercommon.ERROR_NOT_AUTHORISED, http.StatusUnauthorized)
		return
	}
	profile, err := service.profileRepository.FindProfileByUserId(user.ID)
	if err != nil {
		profile := &datamodel.UserProfile{
			UserID: user.ID,
			Name:   user.Name,
		}
		err = service.profileRepository.AddProfile(profile)
		if err != nil {
			slog.Error("can't create profile. Error: " + err.Error())
			http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
			return
		}
	}
	profileData := ProfileData{
		ID:     profile.ID,
		Name:   profile.Name,
		Avatar: profile.Avatar,
	}
	resStr, err := json.Marshal(profileData)
	if err != nil {
		slog.Error("Error. Can't convert projects to json. Error: " + err.Error())
		http.Error(w, servercommon.ERROR_INTERNAL, http.StatusInternalServerError)
		return
	}
	w.Write(resStr)
}
