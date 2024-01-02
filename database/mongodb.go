package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MONGODB_HOST=mongodb
// MONGODB_PORT=27017
// MONGODB_NAME=testdb

type MongoDB[T any] struct {
	Client         *mongo.Client
	DatabaseName   string
	CollectionName string
}

// NewMongoDB takes an env file and returns mongo client
func NewMongoDB[T any](env, collectionName string) (MongoDB[T], error) {
	var db string
	viper.SetConfigFile(env)
	err := viper.ReadInConfig()
	if err != nil {
		return MongoDB[T]{}, err
	}
	host := viper.Get("MONGODB_HOST")
	port := viper.Get("MONGODB_PORT")
	name := viper.Get("MONGODB_NAME")
	val, ok := name.(string)
	if ok {
		db = val
	} else {
		return MongoDB[T]{}, errors.New("could not get database name")
	}

	connectionURI := fmt.Sprintf("mongodb://%s:%s", host, port)

	//set the client options
	clientOptions := options.Client().ApplyURI(connectionURI)

	//connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return MongoDB[T]{}, err
	}

	//test the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return MongoDB[T]{}, err
	}

	//return the client
	return MongoDB[T]{
		Client:         client,
		DatabaseName:   db,
		CollectionName: collectionName,
	}, nil
}

func (m *MongoDB[T]) Insert(ctx context.Context, document T) error {
	collection := m.Client.Database(m.DatabaseName).Collection(m.CollectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB[T]) FindByID(ctx context.Context, id uuid.UUID) (T, error) {
	collection := m.Client.Database(m.DatabaseName).Collection(m.CollectionName)

	var document T
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&document)
	return document, err
}

func (m *MongoDB[T]) Search(ctx context.Context, filter bson.D, opts *options.FindOptions) ([]T, error) {
	results := make([]T, 0)
	if filter == nil {
		return results, errors.New("filter cannot be nil")
	}
	collection := m.Client.Database(m.DatabaseName).Collection(m.CollectionName)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &results)
	if err != nil {
		fmt.Println("WTF", err)
		return results, err
	}

	return results, nil
}

func (m *MongoDB[T]) All(ctx context.Context, opts *options.FindOptions) ([]T, error) {
	results := make([]T, 0)
	filter := bson.D{{}}
	collection := m.Client.Database(m.DatabaseName).Collection(m.CollectionName)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (m *MongoDB[T]) Update(ctx context.Context, id uuid.UUID, document T) (*mongo.UpdateResult, error) {
	collection := m.Client.Database(m.DatabaseName).Collection(m.CollectionName)

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		document,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *MongoDB[T]) Delete(ctx context.Context, filter bson.D, opts *options.DeleteOptions) error {
	collection := m.Client.Database(m.DatabaseName).Collection(m.CollectionName)
	_, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB[T]) Drop(ctx context.Context) error {
	collection := m.Client.Database(m.DatabaseName).Collection(m.CollectionName)
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB[T]) Count(ctx context.Context) (int64, error) {
	return m.Client.Database(m.DatabaseName).Collection(m.CollectionName).EstimatedDocumentCount(ctx)
}
