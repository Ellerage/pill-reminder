package utils

import (
	"context"
	"pill-reminder/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetUserByChatId(db *mongo.Database, chatId int64) (model.User, error) {
	var user model.User

	err := db.Collection("users").FindOne(context.TODO(), bson.M{"chatId": chatId}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
