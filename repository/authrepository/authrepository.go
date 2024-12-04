package authrepository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository struct {
	database *gorm.DB
	notifier servercommon.INotifier
}

func NewAuthRepository(db *gorm.DB, notifier servercommon.INotifier) (AuthRepositorier, error) {
	rep := AuthRepository{
		database: db,
		notifier: notifier,
	}
	err := rep.database.AutoMigrate(&AuthUser{})
	return &rep, err
}

func (repo *AuthRepository) AddUser(user *datamodel.User) error {
	authuser, err := getNewAuthUser(user)
	if err != nil {
		return err
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(authuser)
	if result.Error != nil {
		_, err := repo.FindUserByName(authuser.Name)
		if err == nil {
			return errors.New(servercommon.ERROR_DUPLICATE_USER_NAME)
		}
		_, err = repo.FindUserByEmail(authuser.Email)
		if err == nil {
			return errors.New(servercommon.ERROR_DUPLICATE_EMAIL)
		}
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	err = getDMUser(authuser, user)
	return err
}

func (repo *AuthRepository) FindUser(id uint) (*datamodel.User, error) {
	user := AuthUser{}
	result := repo.database.First(&user, id)
	return getDMByResult(result, &user)
}

func (repo *AuthRepository) FindUserByName(name string) (*datamodel.User, error) {
	user := AuthUser{}
	result := repo.database.First(&user, "name=?", name)
	return getDMByResult(result, &user)
}

func (repo *AuthRepository) FindUserByEmail(email string) (*datamodel.User, error) {
	user := AuthUser{}
	result := repo.database.First(&user, "email=?", email)
	return getDMByResult(result, &user)
}

func (repo *AuthRepository) FindUserByNameEmail(name string) (*datamodel.User, error) {
	user := AuthUser{}
	result := repo.database.Where("name=? OR email=?", name, name).First(&user)
	return getDMByResult(result, &user)
}

func (repo *AuthRepository) DeleteUser(id uint) error {
	result := repo.database.Delete(&AuthUser{}, id)
	return result.Error
}

func (repo *AuthRepository) UpdateAccountStatus(id uint, status datamodel.AccountStatus) error {
	result := repo.database.Model(&AuthUser{}).Set(servercommon.NOTIFIER_KEY, repo.notifier).Where("ID=?", id).Update("acc_status", getStatus(status))
	return result.Error
}
