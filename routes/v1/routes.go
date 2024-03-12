package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/josephmakin/monerohub/handlers"
	"github.com/josephmakin/monerohub/services"
)

func SetupRoutes(router *gin.Engine) {
    paymentsHandler := handlers.NewPaymentsHandler(
        context.TODO(),
        services.Collections["payments"],
        services.RedisClient,
    )

    v1 := router.Group("/api/v1")
    {
        v1.GET("/payment/:id", paymentsHandler.GetOnePaymentHandler)
        v1.POST("/payment", paymentsHandler.CreateOnePaymentHandler)
        v1.POST("/transaction", paymentsHandler.AddOneTransactionHandler)
    }
}
