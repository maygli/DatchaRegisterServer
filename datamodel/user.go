package datamodel

import (
	"datcha/servercommon"
	"net/mail"
	"regexp"
)

type AccountStatus string

const (
	AS_NEED_CONFIRM       AccountStatus = "need_confirm"
	AS_CONFIRMED          AccountStatus = "confirmed"
	AS_RESET_PASSWORD_REQ AccountStatus = "reset_password_req"
	AS_RESET_PASSWORD     AccountStatus = "reset_password"
)

type User struct {
	ID          uint
	Name        string
	Email       string
	Password    string
	AccStatus   AccountStatus
	GoogleId    string
	VkId        string
	TelegrammId string
	YandexId    string
}

func (status AccountStatus) isValid() bool {
	if status != AS_CONFIRMED && status != AS_NEED_CONFIRM &&
		status != AS_RESET_PASSWORD && status != AS_RESET_PASSWORD_REQ {
		return false
	}
	return true
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
	return user.AccStatus.isValid()
}

func (user *User) SetServiceUserId(service string, serviceUserId string) {
	switch service {
	case servercommon.GOOGLE_SERVICE_NAME:
		user.GoogleId = serviceUserId
	case servercommon.VK_SERVICE_NAME:
		user.VkId = serviceUserId
	case servercommon.TELEGRAM_SERVICE_NAME:
		user.TelegrammId = serviceUserId
	case servercommon.YANDEX_SERVICE_NAME:
		user.YandexId = serviceUserId
	}
}

func (user User) HasServiceId(service string) bool {
	switch service {
	case servercommon.GOOGLE_SERVICE_NAME:
		return user.GoogleId != ""
	case servercommon.VK_SERVICE_NAME:
		return user.VkId != ""
	case servercommon.TELEGRAM_SERVICE_NAME:
		return user.TelegrammId != ""
	case servercommon.YANDEX_SERVICE_NAME:
		return user.YandexId != ""
	}
	return false
}
