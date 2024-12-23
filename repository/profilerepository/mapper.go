package profilerepository

import (
	"datcha/datamodel"

	"github.com/jinzhu/copier"
)

func getDMUserProfile(profile *RepProfile) (*datamodel.UserProfile, error) {
	dmProfile := datamodel.UserProfile{}
	err := copier.Copy(&dmProfile, profile)
	return &dmProfile, err
}

func getRepUserProfile(profile *datamodel.UserProfile) (*RepProfile, error) {
	repProfile := RepProfile{}
	err := copier.Copy(&repProfile, profile)
	return &repProfile, err
}

func updateDMProfile(profile *RepProfile, dmProfile *datamodel.UserProfile) error {
	err := copier.Copy(dmProfile, profile)
	return err
}
