package seeds

import (
	"context"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserParams struct {
	TimeToNotify          *string
	Status                *string
	RemindIntervalMinutes *int64
}

func UserSeed(db *mongo.Database, initial UserParams) model.User {
	remindInterval := int64(1)

	if initial.RemindIntervalMinutes != nil {
		remindInterval = *initial.RemindIntervalMinutes
	}

	user := model.User{
		ChatId:         gofakeit.Int64(),
		Timezone:       "UTC",
		TimeToNotify:   gofakeit.Date().Format("15:04"),
		Status:         string(enums.UserStatusIdle),
		RemindInterval: remindInterval,
	}

	if initial.TimeToNotify != nil {
		user.TimeToNotify = *initial.TimeToNotify
	}

	if initial.Status != nil {
		user.Status = *initial.Status
	}

	db.Collection("users").InsertOne(context.TODO(), user)

	return user
}
