package cronnotifier

import (
	"fmt"
	"log/slog"
	"pill-reminder/internal/utils"
)

func (n *NotifierService) reminderFn(chatId int64, repeatCronExp string) func() {
	return func() {
		taken, err := n.deps.PillDayService.IsTakenToday(chatId)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		if !taken {
			n.deps.Notifier.SendMessage(chatId, utils.GetI18nMessage("firstNotification"))

			subID, err := n.cron.AddFunc(repeatCronExp, func() {
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

			if err != nil {
				slog.Error(err.Error())
			}

			slog.Info(fmt.Sprintf("Repeated cron: %s, subID: %d", repeatCronExp, subID))

			n.subIDs[chatId] = subID
		}
	}
}
