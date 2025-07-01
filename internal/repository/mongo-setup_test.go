package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func SetupMongo(t *testing.T) (*mongo.Database, func()) {
	t.Helper()

	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7")

	require.NoError(t, err)

	uri, err := mongodbContainer.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	require.NoError(t, err)

	db := client.Database("testdb")

	teardown := func() {
		_ = client.Disconnect(ctx)
		_ = mongodbContainer.Terminate(ctx)
	}

	return db, teardown
}
