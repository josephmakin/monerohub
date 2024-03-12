package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Collections		map[string]*mongo.Collection
)

func InitMongo(uri, dbName string, colls map[string]string) error {
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
		return err
    }

	Collections = make(map[string]*mongo.Collection)

    db := client.Database(dbName)
	for name, collName := range colls {
		Collections[name] = db.Collection(collName)
	}

	return nil
}
