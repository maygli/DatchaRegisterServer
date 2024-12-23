package profilerepository

import "datcha/datamodel"

type ProfileRepositorier interface {
	AddProfile(profile *datamodel.UserProfile) error
	FindProfileByUserId(userId uint) (*datamodel.UserProfile, error)
	UpdateProfile(profile *datamodel.UserProfile) error
	SetAvatar(profileId uint, avatar string) error
	SetName(profileId uint, name string) error
}
