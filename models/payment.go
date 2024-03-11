package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Address string 	`json:"address" bson:"address"`
	CallbackURL string `json:"callbackurl" bson:"callbackurl"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Transactions []Transaction `json:"transactions" bson:"transactions"`
}
