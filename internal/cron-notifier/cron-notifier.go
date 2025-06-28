package cronnotifier

import (
	"fmt"
	"log"
	"pill-reminder/configs"
	"pill-reminder/internal/service"
	tgbotapi "pill-reminder/internal/tgBotAPI"
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

func ReminderJob(deps NotifierDeps) func() {
	return func() {
		isTakenToday, err := deps.PillDayService.IsTakenToday()

		if err != nil {
			fmt.Println(err)
		}

		if !isTakenToday {
			tgbotapi.SendMessage(deps.Config.MY_CHAT_ID, "Take a pill")

			subCronId, _ = c.AddFunc(repeatedNotificationCron, func() {
				isTakenToday, err := deps.PillDayService.IsTakenToday()

				if err != nil {
					fmt.Println(err)
				}

				if isTakenToday {
					c.Remove(subCronId)
				} else {
					tgbotapi.SendMessage(deps.Config.MY_CHAT_ID, "Reminder take a pill")
				}
			})
		}

	}
}

func RegisterCronNotifier(deps NotifierDeps) {
	loc, _ := time.LoadLocation(timezone)
	c = cron.New(cron.WithLocation(loc))

	ReminderJob(deps)()

	c.AddFunc(firstNotificationCron, ReminderJob(deps))
	c.Start()

	log.Println("Cron started")
}
