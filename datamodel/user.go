package datamodel

import (
	"datcha/servercommon"
	"net/mail"
	"regexp"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountStatus string

const (
	AS_NEED_CONFIRM       AccountStatus = "need_confirm"
	AS_CONFIRMED          AccountStatus = "confirmed"
	AS_RESET_PASSWORD_REQ AccountStatus = "reset_password_req"
	AS_RESET_PASSWORD     AccountStatus = "reset_password"
)

type User struct {
	gorm.Model
	Name      string        `gorm:"unique;not null;type:varchar(100);default:null"`
	Email     string        `gorm:"unique;not null;type:varchar(100);default:null"`
	Password  string        `gorm:"not null;type:varchar(100);default:null"`
	AccStatus AccountStatus `gorm:"not null;type:varchar(100);default:need_confirm"`
	TestField string        `gorm:"varchar(100);default:null"`
}

func (user User) isValid() bool {
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

func (user User) GetHashedPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	notifierInt, hasNot := tx.Get(servercommon.NOTIFIER_KEY)
	if !hasNot {
		return nil
	}
	notifier, ok := notifierInt.(servercommon.INotifier)
	if !ok {
		return nil
	}
	notifier.Notify([]byte("user added" + strconv.FormatUint(uint64(u.ID), 10)))
	return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	notifierInt, hasNot := tx.Get(servercommon.NOTIFIER_KEY)
	if !hasNot {
		return nil
	}
	notifier, ok := notifierInt.(servercommon.INotifier)
	if !ok {
		return nil
	}
	notifier.Notify([]byte("account status changed : " + u.AccStatus))
	return nil
}
