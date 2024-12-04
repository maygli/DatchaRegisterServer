package authrepository

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"datcha/datamodel"
)

func getNewAuthUser(user *datamodel.User) (*AuthUser, error) {
	authuser := AuthUser{}
	err := copier.Copy(&authuser, user)
	return &authuser, err
}

func getNewDMUser(user *AuthUser) (*datamodel.User, error) {
	dmuser := datamodel.User{}
	err := getDMUser(user, &dmuser)
	return &dmuser, err
}

func getDMUser(authuser *AuthUser, user *datamodel.User) error {
	err := copier.Copy(&user, authuser)
	return err
}

func getDMByResult(result *gorm.DB, user *AuthUser) (*datamodel.User, error) {
	if result.Error != nil {
		return &datamodel.User{}, result.Error
	}
	return getNewDMUser(user)
}

func getStatus(status datamodel.AccountStatus) string {
	return string(status)
}
