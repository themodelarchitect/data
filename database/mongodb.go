package database

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
)

// export MONGOGB_CONNECTION_PORT=27017
// export MONGOGB_CONNECTION_HOST=localhost
// export MONGOGB_CONNECTION_DB=testdb

var ErrEnvNotFound = errors.New("environment variable no found")

type MongoDB struct {
	Client   *mongo.Client
	Database string
}

func lookupEnv(name string) (string, error) {
	value, isSet := os.LookupEnv(name)
	if !isSet {
		return "", ErrEnvNotFound
	}
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return "", ErrEnvNotFound
	}
	return value, ErrEnvNotFound
}

func NewMongoDB() (MongoDB, error) {
	host, err := lookupEnv("MONGOGB_CONNECTION_HOST")
	if err != nil {
		return MongoDB{}, err
	}
	port, err := lookupEnv("MONGOGB_CONNECTION_PORT")
	if err != nil {
		return MongoDB{}, err
	}
	db, err := lookupEnv("MONGOGB_CONNECTION_DB")
	if err != nil {
		return MongoDB{}, err
	}

	connectionURI := fmt.Sprintf("mongodb://%s:%s", host, port)

	//set the client options
	clientOptions := options.Client().ApplyURI(connectionURI)

	//connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return MongoDB{}, err
	}

	//test the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return MongoDB{}, err
	}

	//return the client
	return MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

func (m *MongoDB) InsertEntity(ctx context.Context, collectionName string, document interface{}) error {
	collection := m.Client.Database(m.Database).Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) SearchEntity(ctx context.Context, collectionName string, filter interface{}, result interface{}) error {
	collection := m.Client.Database(m.Database).Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, result)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) DeleteEntity(ctx context.Context, collectionName string, filter interface{}) error {
	collection := m.Client.Database(m.Database).Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
