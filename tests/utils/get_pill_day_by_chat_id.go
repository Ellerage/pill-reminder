package utils

import (
	"context"
	"pill-reminder/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetPillDayByChatId(db *mongo.Database, chatId int64) (model.PillDay, error) {
	var pillDay model.PillDay

	nowDate := time.Now().UTC().Format("2006-01-02")

	err := db.Collection("pill-day").FindOne(context.TODO(), bson.M{"chatId": chatId, "date": nowDate}).Decode(&pillDay)
	if err != nil {
		return model.PillDay{}, err
	}

	return pillDay, nil
}
