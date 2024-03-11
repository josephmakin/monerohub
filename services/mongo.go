package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database string = "monerohub"
var PaymentsCollection *mongo.Collection

func InitMongo(uri string) error {
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
		return err
    }

    PaymentsCollection = client.Database(database).Collection("payments")

	return nil
}
