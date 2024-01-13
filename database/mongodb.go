package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// MONGODB_HOST=mongodb
// MONGODB_PORT=27017
// MONGODB_NAME=testdb

type MongoDB[T any] struct {
	Client       *mongo.Client
	DatabaseName string
}

// NewMongoDB takes an env file and returns mongo client
func NewMongoDB[T any]() (MongoDB[T], error) {
	host := os.Getenv("MONGODB_HOST")
	port := os.Getenv("MONGODB_PORT")
	name := os.Getenv("MONGODB_NAME")
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
		Client:       client,
		DatabaseName: name,
	}, nil
}

func (m *MongoDB[T]) Insert(ctx context.Context, collectionName string, document T) error {
	collection := m.Client.Database(m.DatabaseName).Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB[T]) FindByID(ctx context.Context, collectionName string, id uuid.UUID) (T, error) {
	collection := m.Client.Database(m.DatabaseName).Collection(collectionName)

	var document T
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&document)
	return document, err
}

func (m *MongoDB[T]) Search(ctx context.Context, collectionName string, filter bson.D, opts *options.FindOptions) ([]T, error) {
	results := make([]T, 0)
	if filter == nil {
		return results, errors.New("filter cannot be nil")
	}
	collection := m.Client.Database(m.DatabaseName).Collection(collectionName)
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

func (m *MongoDB[T]) All(ctx context.Context, collectionName string, opts *options.FindOptions) ([]T, error) {
	results := make([]T, 0)
	filter := bson.D{{}}
	collection := m.Client.Database(m.DatabaseName).Collection(collectionName)
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

func (m *MongoDB[T]) Update(ctx context.Context, collectionName string, id uuid.UUID, document T) (*mongo.UpdateResult, error) {
	collection := m.Client.Database(m.DatabaseName).Collection(collectionName)

	filter := bson.D{{"_id", id}}
	//update := bson.D{{"$set", bson.D{{"email", "newemail@example.com"}}}}
	update := bson.D{{"$set", document}}
	result, err := collection.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *MongoDB[T]) Delete(ctx context.Context, collectionName string, filter bson.D, opts *options.DeleteOptions) error {
	collection := m.Client.Database(m.DatabaseName).Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB[T]) Drop(ctx context.Context, collectionName string) error {
	collection := m.Client.Database(m.DatabaseName).Collection(collectionName)
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB[T]) Count(ctx context.Context, collectionName string) (int64, error) {
	return m.Client.Database(m.DatabaseName).Collection(collectionName).EstimatedDocumentCount(ctx)
}
