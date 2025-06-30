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

const repeatedNotificationCron = "*/20 * * * *"

type CronNotifier struct {
	cron     *cron.Cron
	cronIDs  map[int64]cron.EntryID
	subIDs   map[int64]cron.EntryID
	deps     NotifierDeps
	timezone *time.Location
}

type NotifierDeps struct {
	PillDayService *service.PillDayService
	UserService    *service.UserService
	Config         *configs.Config
}

func NewCronNotifier(deps NotifierDeps) (*CronNotifier, error) {
	loc, err := time.LoadLocation(deps.Config.TIMEZONE)
	if err != nil {
		return nil, err
	}

	return &CronNotifier{
		cron:     cron.New(cron.WithLocation(loc)),
		cronIDs:  make(map[int64]cron.EntryID),
		subIDs:   make(map[int64]cron.EntryID),
		deps:     deps,
		timezone: loc,
	}, nil
}

func (n *CronNotifier) sendRepeatedReminder(chatId int64) {
	tgbotapi.SendMessage(chatId, utils.GetI18nMessage("firstNotification"))

	subID, _ := n.cron.AddFunc(repeatedNotificationCron, func() {
		isTakenToday, err := n.deps.PillDayService.IsTakenToday(chatId)

		if err != nil {
			slog.Error(err.Error())
			return
		}
		if isTakenToday {
			n.cron.Remove(n.subIDs[chatId])
		} else {
			tgbotapi.SendMessage(chatId, utils.GetI18nMessage("reminderNotification"))
		}
	})

	n.subIDs[chatId] = subID
}

func (n *CronNotifier) reminderFn(chatId int64) func() {
	return func() {
		taken, err := n.deps.PillDayService.IsTakenToday(chatId)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		if !taken {
			n.sendRepeatedReminder(chatId)
		}
	}
}

func (n *CronNotifier) AddOrUpdateCron(chatId int64, time time.Time) (cron.EntryID, string, error) {
	if oldID, ok := n.cronIDs[chatId]; ok {
		n.cron.Remove(oldID)
	}

	cronStr := n.GetCronExpFromTime(time)
	id, err := n.cron.AddFunc(cronStr, n.reminderFn(chatId))
	if err != nil {
		return 0, "", err
	}

	n.cronIDs[chatId] = id

	return id, cronStr, nil
}

func (n *CronNotifier) GetCronExpFromTime(t time.Time) string {
	return fmt.Sprintf("%d %d * * *", t.Minute(), t.Hour())
}

func (n *CronNotifier) Start() {
	users, err := n.deps.UserService.GetAll()
	if err != nil {
		log.Println("Failed to fetch users:", err)
		return
	}

	for _, user := range users {
		timeToNotify, err := utils.ConvertTimeToTbilisi(user.TimeToNotify, user.Timezone)
		if err != nil {
			slog.Error("Error creating cron", "error", err)
			continue
		}

		_, cronExp, err := n.AddOrUpdateCron(user.ChatId, timeToNotify)
		if err != nil {
			slog.Error("Failed to add cron", "chatId", user.ChatId, "error", err)
			continue
		}

		slog.Info("Created cron", "chatId", user.ChatId, "cron", cronExp)
	}

	n.cron.Start()
	log.Println("Cron started")
}
