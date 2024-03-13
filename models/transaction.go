package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Transaction struct {
    Address string `json:"address,omitempty" bson:"-"`
    Amount float32 `json:"amount" bson:"amount"`
    Timestamp time.Time `json:"timestamp" bson:"timestamp"`
    TxID string `json:"txid" bson:"txid"`
}

func (t *Transaction) Relay(callbackURL string) error {
    requestBody, err := json.Marshal(t)
    if err != nil {
        return err
    }

    _, err = http.Post(callbackURL, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return err
    }
    return nil
}
