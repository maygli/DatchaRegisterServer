package authrepository

import "datcha/datamodel"

type AuthRepositorier interface {
	AddUser(user *datamodel.User) error
	FindUser(id uint) (*datamodel.User, error)
	FindUserByName(name string) (*datamodel.User, error)
	FindUserByEmail(email string) (*datamodel.User, error)
	FindUserByNameEmail(name string) (*datamodel.User, error)
	DeleteUser(id uint) error
	UpdateAccountStatus(id uint, status datamodel.AccountStatus) error
}
