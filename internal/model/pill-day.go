package model

type PillDay struct {
	Date         string  `bson:"date"`
	TimeOfTaking *string `bson:"timeOfTaking,omitempty"`
	ChatId       int64   `bson:"chatId"`
}

func (p *PillDay) HasTimeOfTaking() bool {
	return p.TimeOfTaking != nil && *p.TimeOfTaking != ""
}
