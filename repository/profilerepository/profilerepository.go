package profilerepository

import (
	"datcha/datamodel"
	"datcha/servercommon"
	"errors"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	database *gorm.DB
	notifier servercommon.INotifier
}

func NewProfileRepository(db *gorm.DB, notifier servercommon.INotifier) (ProfileRepositorier, error) {
	rep := ProfileRepository{
		database: db,
		notifier: notifier,
	}
	err := rep.database.AutoMigrate(&RepProfile{})
	return &rep, err
}

func (repo *ProfileRepository) AddProfile(profile *datamodel.UserProfile) error {
	repProfile, err := getRepUserProfile(profile)
	if err != nil {
		return err
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Create(repProfile)
	if result.Error != nil {
		return errors.New(servercommon.ERROR_INTERNAL)
	}
	err = updateDMProfile(repProfile, profile)
	return err
}

func (repo *ProfileRepository) FindProfileByUserId(userId uint) (*datamodel.UserProfile, error) {
	profile := RepProfile{}
	result := repo.database.Where("user_id = ?", userId).First(&profile)
	if result.Error != nil {
		return &datamodel.UserProfile{}, result.Error
	}
	dmProfile, err := getDMUserProfile(&profile)
	return dmProfile, err
}

func (repo *ProfileRepository) UpdateProfile(profile *datamodel.UserProfile) error {
	repProfile, err := getRepUserProfile(profile)
	if err != nil {
		return err
	}
	result := repo.database.Set(servercommon.NOTIFIER_KEY, repo.notifier).Save(repProfile)
	return result.Error
}

func (repo *ProfileRepository) SetAvatar(profileId uint, avatar string) error {
	result := repo.database.Model(&RepProfile{ID: profileId}).Update("avatar", avatar)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (repo *ProfileRepository) SetName(profileId uint, name string) error {
	result := repo.database.Model(&RepProfile{ID: profileId}).Update("name", name)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
