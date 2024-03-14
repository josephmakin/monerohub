package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/josephmakin/monerohub/models"
	"github.com/josephmakin/monerohub/services"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentsHandler struct {
    collection *mongo.Collection
    ctx context.Context
    redisClient *redis.Client
}

func NewPaymentsHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *PaymentsHandler {
    return &PaymentsHandler{
        collection: collection,
        ctx: ctx,
        redisClient: redisClient,
    }
}

func (paymentsHandler *PaymentsHandler) GetOnePaymentHandler (c *gin.Context) {
    id := c.Param("id")
    
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": err.Error(),
        })
        return
    }

    var payment models.Payment
    err = paymentsHandler.collection.FindOne(paymentsHandler.ctx, bson.M{
        "_id": objectID,
    }).Decode(&payment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, payment)
}

func (paymentsHandler *PaymentsHandler) CreateOnePaymentHandler (c *gin.Context) {
    var payment models.Payment

    if err := c.ShouldBindJSON(&payment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    resp, err := services.Monero.CreateAddress(&wallet.RequestCreateAddress{
    	AccountIndex: 0,
    	Label:        "",
    })
    payment.Address = resp.Address
    payment.ID = primitive.NewObjectID()
    payment.Timestamp = time.Now()
    payment.Transactions = []models.Transaction{}

    _, err = paymentsHandler.collection.InsertOne(paymentsHandler.ctx, payment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    log.Println("Remove data from Redis")
    paymentsHandler.redisClient.Del("payments")

    c.JSON(http.StatusOK, payment)
}

func (handler *PaymentsHandler) AddOneTransactionHandler (c *gin.Context) {
    var transaction models.Transaction

    if err := c.ShouldBindJSON(&transaction); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    transactionDetails, err := services.Monero.GetTransferByTxID(&wallet.RequestGetTransferByTxID{TxID: transaction.TxID})
    if err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": err.Error(),
        })
        return
    }

    transaction = models.Transaction{
        Address: transactionDetails.Transfer.Address,
        Amount: float32(transactionDetails.Transfer.Amount),
        Timestamp: time.Unix(int64(transactionDetails.Transfer.Timestamp), 0),
        TxID: transactionDetails.Transfer.TxID,
    }

    filter := bson.M{
        "address": transaction.Address,
        "transactions.txid": bson.M{"$ne": transaction.TxID},
    }

    update := bson.M{
        "$addToSet": bson.M{
            "transactions": transaction,
        },
    }

    opts := options.Update().SetUpsert(true)

    _, err = handler.collection.UpdateOne(handler.ctx, filter, update, opts)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
    }

    var payment models.Payment
    err = handler.collection.FindOne(handler.ctx, bson.M{
        "address": transaction.Address,
    }).Decode(&payment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    err = transaction.Relay(payment.CallbackURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, transaction)
}
