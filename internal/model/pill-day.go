package model

type PillDay struct {
	Date         string  `bson:"date,omitempty"`
	TimeOfTaking *string `bson:"timeOfTaking,omitempty"`
}

func (p *PillDay) HasTimeOfTaking() bool {
	return p.TimeOfTaking != nil && *p.TimeOfTaking != ""
}
