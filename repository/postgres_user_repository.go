package repository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"

	_ "github.com/lib/pq"
)

func (repo *DatchaSqlRepository) AddUser(user *datamodel.User) error {
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(user)
	if result.Error != nil {
		_, err := repo.FindUserByName(user.Name)
		if err == nil {
			return errors.New(servercommon.ERROR_DUPLICATE_USER_NAME)
		}
		_, err = repo.FindUserByEmail(user.Email)
		if err == nil {
			return errors.New(servercommon.ERROR_DUPLICATE_EMAIL)
		}
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	return result.Error
}

func (repo *DatchaSqlRepository) FindUser(id uint) (datamodel.User, error) {
	user := datamodel.User{}
	result := repo.database.First(&user, id)
	return user, result.Error
}

func (repo *DatchaSqlRepository) FindUserByName(name string) (datamodel.User, error) {
	user := datamodel.User{}
	ctx := repo.database.First(&user, "name=?", name)
	return user, ctx.Error
}

func (repo *DatchaSqlRepository) FindUserByEmail(email string) (datamodel.User, error) {
	user := datamodel.User{}
	ctx := repo.database.First(&user, "email=?", email)
	return user, ctx.Error
}

func (repo *DatchaSqlRepository) FindUserByNameEmail(name string) (datamodel.User, error) {
	user := datamodel.User{}
	ctx := repo.database.Where("name=? OR email=?", name, name).First(&user)
	return user, ctx.Error
}

func (repo *DatchaSqlRepository) DeleteUser(id uint) error {
	ctx := repo.database.Delete(&datamodel.User{}, id)
	return ctx.Error
}

func (repo *DatchaSqlRepository) UpdateAccountStatus(id uint, status datamodel.AccountStatus) error {
	ctx := repo.database.Model(&datamodel.User{}).Set(servercommon.NOTIFIER_KEY, repo.notifier).Where("ID=?", id).Update("acc_status", status)
	return ctx.Error
}
