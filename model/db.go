package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	db *mongo.Database
}

func NewConnection(cs string) (*Connection, error) {
	clientOptions := options.Client().ApplyURI(cs)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to MongoDB!")

	db := client.Database("scanner")

	return &Connection{db}, nil
}

func (c *Connection) GetCollection() *mongo.Collection {
	return c.db.Collection("program")
}
