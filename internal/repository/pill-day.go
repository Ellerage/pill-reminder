package repository

import (
	"context"
	"log/slog"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PillDayRepo struct {
	db *mongo.Database
}

func NewPillDayRepo(db *mongo.Database) *PillDayRepo {
	return &PillDayRepo{db: db}
}

func (repo *PillDayRepo) GetByDateAndChatId(chatId int64, date time.Time) (*model.PillDay, error) {
	var result *model.PillDay

	formattedDate := date.Format("2006-01-02")

	err := repo.db.Collection("pill-day").FindOne(context.TODO(), bson.M{"date": formattedDate, "chatId": chatId}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *PillDayRepo) Create(chatId int64, timeOfTaking *time.Time) error {
	var formattedTime *string

	if timeOfTaking != nil {
		str := timeOfTaking.Format("15:04")
		formattedTime = &str
	} else {
		formattedTime = nil
	}

	pillDay := model.PillDay{Date: utils.GetFormattedNowDate(nil), TimeOfTaking: formattedTime, ChatId: chatId}

	_, err := repo.db.Collection("pill-day").InsertOne(context.TODO(), pillDay)

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (repo *PillDayRepo) UpdateTimeByDate(chatId int64, dateTime time.Time) error {
	toUpdate := bson.M{"$set": bson.M{"timeOfTaking": dateTime.Format("15:04")}}

	_, err := repo.db.Collection("pill-day").UpdateOne(context.TODO(), bson.M{"date": dateTime.Format("2006-01-02"), "chatId": chatId}, toUpdate)

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
