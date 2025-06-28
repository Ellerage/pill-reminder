package repository

import (
	"context"
	"fmt"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"

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

func (repo *PillDayRepo) GetByDate(date string) (*model.PillDay, error) {
	var result *model.PillDay

	// TODO: add database name to config?
	err := repo.db.Collection("pill-day").FindOne(context.TODO(), bson.M{"date": date}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *PillDayRepo) Create(timeOfTaking *string) error {
	pillDay := model.PillDay{Date: utils.GetNowDateTbilisi(), TimeOfTaking: timeOfTaking}

	// TODO: add database name to config?
	_, err := repo.db.Collection("pill-day").InsertOne(context.TODO(), pillDay)

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (repo *PillDayRepo) UpdateTimeByDate(date string, time string) error {
	toUpdate := bson.M{"$set": bson.M{"timeOfTaking": time}}

	// TODO: add database name to config?
	_, err := repo.db.Collection("pill-day").UpdateOne(context.TODO(), bson.M{"date": date}, toUpdate)

	if err != nil {
		fmt.Println(err)
	}

	return err
}
