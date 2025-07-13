package repository

import (
	"context"
	"log/slog"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils/enums"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) GetAll() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := repo.db.Collection("users").Find(ctx, bson.M{
		"status": bson.M{
			"$ne": string(enums.UserStatusInactive),
		},
	})

	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepo) GetByChatId(chatId int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user *model.User

	err := repo.db.Collection("users").FindOne(ctx, bson.M{"chatId": chatId}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepo) Create(toCreate model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.db.Collection("users").InsertOne(ctx, toCreate)

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (repo *UserRepo) Update(chatId int64, toUpdate model.UserUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	set := bson.M{}

	if toUpdate.Status != nil {
		set["status"] = *toUpdate.Status
	}

	if toUpdate.Timezone != nil {
		set["timezone"] = *toUpdate.Timezone
	}

	if toUpdate.TimeToNotify != nil {
		set["timeToNotify"] = *toUpdate.TimeToNotify
	}

	if toUpdate.RemindInterval != nil {
		set["remindInterval"] = *toUpdate.RemindInterval
	}

	_, err := repo.db.Collection("users").UpdateOne(ctx, bson.M{"chatId": chatId}, bson.M{"$set": set})

	if err != nil {
		slog.Error(err.Error())

		return err
	}

	return nil
}
