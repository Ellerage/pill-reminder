package cronnotifier

import (
	"pill-reminder/internal/utils/enums"
	"time"

	"github.com/robfig/cron/v3"
)

type Notifier interface {
	SendMessage(chatID int64, message string, buttons *enums.SendMessageButtons)
}

type PillDayService interface {
	IsTakenToday(chatID int64) (bool, error)
}

type NotifierService struct {
	cron    *cron.Cron
	cronIDs map[int64]cron.EntryID
	subIDs  map[int64]cron.EntryID
	deps    NotifierParams
}

type NotifierParams struct {
	PillDayService PillDayService
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
