package repository

import (
	"context"
	"log/slog"
	"pill-reminder/internal/model"

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
	cursor, err := repo.db.Collection("users").Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	var users []model.User

	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepo) Create(toCreate model.User) error {
	_, err := repo.db.Collection("users").InsertOne(context.TODO(), toCreate)

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
