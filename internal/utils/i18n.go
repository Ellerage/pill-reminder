package utils

var i18nTexts = map[string]string{
	"firstNotification":    "Take a pill",
	"reminderNotification": "Reminder take a pill",
	"markTakenBtn":         "Take",
}

func GetI18nMessage(tag string) string {
	return i18nTexts[tag]
}
