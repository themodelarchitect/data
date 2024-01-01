package database

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// export MONGODB_HOST=mongodb
// export MONGODB_PORT=27017
// export MONGODB_NAME=testdb

var ErrEnvNotFound = errors.New("environment variable not found")

type MongoDB struct {
	Client   *mongo.Client
	Database string
}

// NewMongoDB takes an env file and returns mongo client
func NewMongoDB(env string) (MongoDB, error) {
	var db string
	viper.SetConfigFile(env)
	err := viper.ReadInConfig()
	if err != nil {
		return MongoDB{}, err
	}
	host := viper.Get("MONGODB_HOST")
	port := viper.Get("MONGODB_PORT")
	name := viper.Get("MONGODB_NAME")
	val, ok := name.(string)
	if ok {
		db = val
	} else {
		return MongoDB{}, errors.New("could not get database name")
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

func (m *MongoDB) Insert(ctx context.Context, collectionName string, document any) error {
	collection := m.Client.Database(m.Database).Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Search(ctx context.Context, collectionName string, filter any, results any) error {
	if filter == nil {
		return errors.New("filter cannot be nil")
	}
	collection := m.Client.Database(m.Database).Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, results)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) All(ctx context.Context, collectionName string, results any) error {
	filter := bson.D{{}}
	collection := m.Client.Database(m.Database).Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, results)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Delete(ctx context.Context, collectionName string, filter any) error {
	collection := m.Client.Database(m.Database).Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Drop(ctx context.Context, collectionName string) error {
	collection := m.Client.Database(m.Database).Collection(collectionName)
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Count(ctx context.Context, collectionName string) (int64, error) {
	return m.Client.Database(m.Database).Collection(collectionName).EstimatedDocumentCount(ctx)
}
