package cronnotifier

import (
	"fmt"
	"log"
	"log/slog"
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

func sendRepeatedReminder(deps NotifierDeps, chatId int64) {
	tgbotapi.SendMessage(chatId, utils.GetI18nMessage("firstNotification"))

	subCronId, _ = c.AddFunc(repeatedNotificationCron, func() {
		isTakenToday, err := deps.PillDayService.IsTakenToday(chatId)

		if err != nil {
			slog.Error(err.Error())
		}

		if isTakenToday {
			c.Remove(subCronId)
		} else {
			tgbotapi.SendMessage(chatId, utils.GetI18nMessage("reminderNotification"))
		}
	})
}

func ReminderNotification(deps NotifierDeps, chatId int64) func() {
	return func() {
		if taken, _ := deps.PillDayService.IsTakenToday(chatId); !taken {
			sendRepeatedReminder(deps, chatId)
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
			slog.Error("Error creating cron", "error", err)
		}

		cronStr := fmt.Sprintf("%d %d * * *", timeToNotify.Minute(), timeToNotify.Hour())

		c.AddFunc(cronStr, ReminderNotification(deps, user.ChatId))

		slog.Info(fmt.Sprintf("Created Cron for chatId %d, at time %s", user.ChatId, cronStr))
	}

	c.Start()

	log.Println("Cron started")
}
