package cronnotifier

import (
	"fmt"
	"log"
	"pill-reminder/configs"
	"pill-reminder/internal/service"
	tgbotapi "pill-reminder/internal/tgBotAPI"
	"pill-reminder/internal/utils"
	"time"

	"github.com/robfig/cron/v3"
)

type NotifierDeps struct {
	PillDayService *service.PillDayService
	UserService    *service.UserService
	Config         *configs.Config
}

var repeatedNotificationCron = "*/20 * * * *"

var c *cron.Cron
var subCronId cron.EntryID

func sendRepeatedReminder(deps NotifierDeps) {
	tgbotapi.SendMessage(deps.Config.MY_CHAT_ID, utils.GetI18nMessage("firstNotification"))

	subCronId, _ = c.AddFunc(repeatedNotificationCron, func() {
		isTakenToday, err := deps.PillDayService.IsTakenToday()

		if err != nil {
			fmt.Println(err)
		}

		if isTakenToday {
			c.Remove(subCronId)
		} else {
			tgbotapi.SendMessage(deps.Config.MY_CHAT_ID, utils.GetI18nMessage("reminderNotification"))
		}
	})
}

func ReminderNotification(deps NotifierDeps) func() {
	return func() {
		if taken, _ := deps.PillDayService.IsTakenToday(); !taken {
			sendRepeatedReminder(deps)
		}
	}
}

func RegisterCronNotifier(deps NotifierDeps) {

	users, err := deps.UserService.GetAll()

	if err != nil {
		log.Println(err)
	}

	loc, _ := time.LoadLocation(deps.Config.TIMEZONE)
	c = cron.New(cron.WithLocation(loc))

	for _, user := range users {
		timeToNotify, err := utils.ConvertTimeToTbilisi(user.TimeToNotify, user.Timezone)

		if err != nil {
			log.Println(err)
		}

		cronStr := fmt.Sprintf("%d %d * * *", timeToNotify.Minute(), timeToNotify.Hour())

		// TODO: extend collection with userId of history
		c.AddFunc(cronStr, ReminderNotification(deps))
	}

	c.Start()

	log.Println("Cron started")
}
