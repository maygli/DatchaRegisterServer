package repositorycommon

import (
	"datcha/servercommon"

	"gorm.io/gorm"
)

const (
	NOTIFIER_KEY string = "notifier"
)

func SendNotify(db *gorm.DB, message string) {
	notifierInt, hasNot := db.Get(NOTIFIER_KEY)
	if !hasNot {
		return
	}
	notifier, ok := notifierInt.(servercommon.INotifier)
	if !ok {
		return
	}
	notifier.Notify([]byte(message))
	return
}
