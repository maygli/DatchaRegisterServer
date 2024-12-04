package authrepository

import (
	"datcha/repository/repositorycommon"
	"fmt"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUser struct {
	gorm.Model
	Name      string `gorm:"unique;not null;type:varchar(100);default:null"`
	Email     string `gorm:"unique;not null;type:varchar(100);default:null"`
	Password  string `gorm:"not null;type:varchar(100);default:null"`
	AccStatus string `gorm:"not null;type:varchar(100);default:need_confirm"`
}

func (user AuthUser) isValid() bool {
	if user.Name == "" && user.Email == "" {
		return false
	}
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return false
	}
	checkNameRe := regexp.MustCompile(`^[_a-z]\\w*$`)
	res := checkNameRe.MatchString(user.Name)
	if !res {
		return false
	}
	return true
}

func (user AuthUser) GetHashedPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
}

func (u *AuthUser) AfterCreate(tx *gorm.DB) error {
	repositorycommon.SendNotify(tx, fmt.Sprintf("used added. Id=%d", u.ID))
	return nil
}

func (u *AuthUser) AfterUpdate(tx *gorm.DB) (err error) {
	repositorycommon.SendNotify(tx, fmt.Sprintf("user account updated. Status: %s", u.AccStatus))
	return nil
}
