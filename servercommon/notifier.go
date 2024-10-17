package servercommon

import "gorm.io/gorm"

func SendNotify(db *gorm.DB, message string) {
	notifierInt, hasNot := db.Get(NOTIFIER_KEY)
	if !hasNot {
		return
	}
	notifier, ok := notifierInt.(INotifier)
	if !ok {
		return
	}
	notifier.Notify([]byte(message))
	return
}
