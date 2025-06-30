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
	SendMessage(chatID int64, message string)
}

type NotifierService struct {
	cron    *cron.Cron
	cronIDs map[int64]cron.EntryID
	subIDs  map[int64]cron.EntryID
	deps    NotifierParams
}

type NotifierParams struct {
	PillDayService *service.PillDayService
	Timezone       string
	Notifier       Notifier
}

func NewCronNotifier(params NotifierParams) (*NotifierService, error) {
	loc, err := time.LoadLocation(params.Timezone)

	if err != nil {
		return nil, err
	}

	return &NotifierService{
		cron:    cron.New(cron.WithLocation(loc)),
		cronIDs: make(map[int64]cron.EntryID),
		subIDs:  make(map[int64]cron.EntryID),
		deps:    params,
	}, nil
}

func (n *NotifierService) SetNotifier(notifier Notifier) {
	n.deps.Notifier = notifier
}

func (n *NotifierService) sendRepeatedReminder(chatId int64) {
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

func (n *NotifierService) reminderFn(chatId int64) func() {
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

func (n *NotifierService) AddOrUpdateCron(chatId int64, time time.Time) error {
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

func (n *NotifierService) GetCronExpFromTime(t time.Time) string {
	return fmt.Sprintf("%d %d * * *", t.Minute(), t.Hour())
}

func (n *NotifierService) Start(users []model.User) {
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
