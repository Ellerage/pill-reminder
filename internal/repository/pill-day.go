package repository

import (
	"context"
	"fmt"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Options struct {
	TimeOfTaking bool
}

type PillDayRepo struct {
	db *mongo.Database
}

func NewPillDayRepo(db *mongo.Database) *PillDayRepo {
	return &PillDayRepo{db: db}
}

func (repo *PillDayRepo) GetByDate(date time.Time) (*model.PillDay, error) {
	var result *model.PillDay

	formattedDate := date.Format("2006-01-02")

	err := repo.db.Collection("pill-day").FindOne(context.TODO(), bson.M{"date": formattedDate}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *PillDayRepo) Create(timeOfTaking *time.Time) error {
	formattedTime := timeOfTaking.Format("15:04")

	pillDay := model.PillDay{Date: utils.GetFormattedNowDate(nil), TimeOfTaking: &formattedTime}

	_, err := repo.db.Collection("pill-day").InsertOne(context.TODO(), pillDay)

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (repo *PillDayRepo) UpdateTimeByDate(date time.Time, time time.Time) error {
	toUpdate := bson.M{"$set": bson.M{"timeOfTaking": time.Format("15:04")}}

	_, err := repo.db.Collection("pill-day").UpdateOne(context.TODO(), bson.M{"date": date.Format("2006-01-02")}, toUpdate)

	if err != nil {
		fmt.Println(err)
	}

	return err
}
