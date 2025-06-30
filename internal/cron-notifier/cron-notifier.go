package cronnotifier

import (
	"fmt"
	"log"
	"log/slog"
	"pill-reminder/internal/model"
	"pill-reminder/internal/service"
	"pill-reminder/internal/utils"
	"time"

	"github.com/robfig/cron/v3"
)

const repeatedNotificationCron = "*/20 * * * *"

type Notifier interface {
	SendMessage(chatId int64, message string)
}

type CronNotifier struct {
	cron    *cron.Cron
	cronIDs map[int64]cron.EntryID
	subIDs  map[int64]cron.EntryID
	deps    NotifierDeps
}

type NotifierDeps struct {
	PillDayService *service.PillDayService
	Timezone       string
	Notifier       Notifier
}

func NewCronNotifier(deps NotifierDeps) (*CronNotifier, error) {
	loc, err := time.LoadLocation(deps.Timezone)

	if err != nil {
		return nil, err
	}

	return &CronNotifier{
		cron:    cron.New(cron.WithLocation(loc)),
		cronIDs: make(map[int64]cron.EntryID),
		subIDs:  make(map[int64]cron.EntryID),
		deps:    deps,
	}, nil
}

func (n *CronNotifier) sendRepeatedReminder(chatId int64) {
	n.deps.Notifier.SendMessage(chatId, utils.GetI18nMessage("firstNotification"))

	subID, _ := n.cron.AddFunc(repeatedNotificationCron, func() {
		isTakenToday, err := n.deps.PillDayService.IsTakenToday(chatId)

		if err != nil {
			slog.Error(err.Error())
			return
		}
		if isTakenToday {
			n.cron.Remove(n.subIDs[chatId])
		} else {
			n.deps.Notifier.SendMessage(chatId, utils.GetI18nMessage("reminderNotification"))
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

func (n *CronNotifier) AddOrUpdateCron(chatId int64, time time.Time) error {
	if oldID, ok := n.cronIDs[chatId]; ok {
		n.cron.Remove(oldID)
		slog.Info(fmt.Sprintf("Cron with ID: %d was removed", oldID))
	}

	cronStr := n.GetCronExpFromTime(time)
	id, err := n.cron.AddFunc(cronStr, n.reminderFn(chatId))
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Cron new added ID: %d, cron exp: %s", id, cronStr))

	n.cronIDs[chatId] = id

	return nil
}

func (n *CronNotifier) GetCronExpFromTime(t time.Time) string {
	return fmt.Sprintf("%d %d * * *", t.Minute(), t.Hour())
}

func (n *CronNotifier) Start(users []model.User) {
	for _, user := range users {
		timeToNotify, err := utils.ConvertTimeToTbilisi(user.TimeToNotify, user.Timezone)
		if err != nil {
			slog.Error("Error creating cron", "error", err)
			continue
		}

		if err := n.AddOrUpdateCron(user.ChatId, timeToNotify); err != nil {
			slog.Error("Failed to add cron", "chatId", user.ChatId, "error", err)
			continue
		}

	}

	n.cron.Start()
	log.Println("Cron started")
}
