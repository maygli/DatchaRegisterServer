package authrepository

import (
	"datcha/datamodel"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	TEST_ID         uint   = 345
	TEST_NAME       string = "name"
	TEST_PASSWORD   string = "password"
	TEST_EMAIL      string = "mail@test.com"
	TEST_ACC_STATUS string = "reset_password_req"
)

func validateDMUser(t *testing.T, user *datamodel.User) {
	assert.Equal(t, user.ID, TEST_ID)
	assert.Equal(t, user.Name, TEST_NAME)
	assert.Equal(t, user.Email, TEST_EMAIL)
	assert.Equal(t, user.Password, TEST_PASSWORD)
	assert.Equal(t, string(user.AccStatus), TEST_ACC_STATUS)
}

func TestConvertToDMUser(t *testing.T) {
	user := AuthUser{
		Model: gorm.Model{
			ID: TEST_ID,
		},
		Name:      TEST_NAME,
		Email:     TEST_EMAIL,
		Password:  TEST_PASSWORD,
		AccStatus: TEST_ACC_STATUS,
	}
	dmUser, err := getNewDMUser(&user)
	assert.NotNil(t, err)
	assert.NotNil(t, dmUser)
	validateDMUser(t, dmUser)
}
