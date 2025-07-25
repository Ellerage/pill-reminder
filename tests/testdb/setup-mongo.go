package testsdb

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Database

func SetupMongo() func() {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		panic(err)
	}

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("testdb")

	teardown := func() {
		_ = client.Disconnect(ctx)
		_ = mongodbContainer.Terminate(ctx)
	}

	MongoClient = db

	return teardown
}

func ClearMongo(db *mongo.Database) {
	collections, _ := db.ListCollectionNames(context.TODO(), bson.D{})
	for _, name := range collections {
		_ = db.Collection(name).Drop(context.TODO())
	}
}
