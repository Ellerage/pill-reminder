package db

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type ConnectMongoOptions struct {
	Uri      string
	DBName   string
	UserName string
	Password string
}

func ConnectMongo(connectOptions ConnectMongoOptions) *mongo.Database {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectOptions.Uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	} else {
		slog.Info("Successfully connected to MongoDB!")
	}

	return client.Database(connectOptions.DBName)
}
