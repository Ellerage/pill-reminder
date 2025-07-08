package cronnotifier

import (
	"fmt"
	"log"
	"log/slog"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
)

func (n *NotifierService) AddOrUpdateCron(chatId int64, timeStr string, repeatCronExp string) error {
	time := utils.GetTimeFromString(timeStr)

	if oldID, ok := n.cronIDs[chatId]; ok {
		n.cron.Remove(oldID)
		slog.Info(fmt.Sprintf("Cron with ID: %d was removed", oldID))
	}

	cronStr := fmt.Sprintf("%d %d * * *", time.Minute(), time.Hour())
	id, err := n.cron.AddFunc(cronStr, n.reminderFn(chatId, repeatCronExp))
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Cron new added ID: %d, cron exp: %s. Repeated sub Cron: %s", id, cronStr, repeatCronExp))

	n.cronIDs[chatId] = id

	return nil
}

func (n *NotifierService) Start(users []model.User) {
	for _, user := range users {
		if err := n.AddOrUpdateCron(user.ChatId, user.TimeToNotify, user.RemindInterval); err != nil {
			slog.Error("Failed to add cron", "chatId", user.ChatId, "error", err)
			continue
		}

	}

	n.cron.Start()
	log.Println("Cron started")
}
