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
	Config         *configs.Config
}

var (
	firstNotificationCron    = "0 17 * * *"
	repeatedNotificationCron = "*/20 * * * *"
	timezone                 = "Asia/Tbilisi"
)

var c *cron.Cron
var subCronId cron.EntryID

func ReminderNotification(deps NotifierDeps) func() {
	return func() {
		isTakenToday, err := deps.PillDayService.IsTakenToday()

		if err != nil {
			fmt.Println(err)
		}

		if !isTakenToday {
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

	}
}

func RegisterCronNotifier(deps NotifierDeps) {
	loc, _ := time.LoadLocation(timezone)
	c = cron.New(cron.WithLocation(loc))

	// TODO: Replace with initial check
	ReminderNotification(deps)()

	c.AddFunc(firstNotificationCron, ReminderNotification(deps))
	c.Start()

	log.Println("Cron started")
}
