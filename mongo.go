package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func connectToMongoDB(uri, dbName, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	cli, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	err = cli.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("error pinging MongoDB: %v", err)
	}

	client = cli
	collection = client.Database(dbName).Collection(collectionName)

	fmt.Println("Connected to MongoDB!")
	return nil
}

func insertFormData(name, phoneno, address, preferredtime string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, bson.M{
		"name":          name,
		"phoneno":       phoneno,
		"address":       address,
		"preferredtime": preferredtime,
	})
	if err != nil {
		return fmt.Errorf("error inserting document: %v", err)
	}

	fmt.Println("Inserted document into MongoDB!")
	return nil
}
