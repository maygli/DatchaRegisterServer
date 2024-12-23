package profileservice

import "datcha/datamodel"

func mergeProfileData(profile *datamodel.UserProfile, data ProfileData) {
	if data.Name != "" {
		profile.Name = data.Name
	}
	if data.Avatar != "" {
		profile.Avatar = data.Avatar
	}
}
