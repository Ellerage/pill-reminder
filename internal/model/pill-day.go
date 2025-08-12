package model

type PillDay struct {
	Date         string  `db:"date"`
	TimeOfTaking *string `db:"timeOfTaking"`
	ChatId       int64   `db:"chatId"`
}

func (p *PillDay) HasTimeOfTaking() bool {
	return p.TimeOfTaking != nil && *p.TimeOfTaking != ""
}
