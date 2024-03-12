package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	v1 "github.com/josephmakin/monerohub/routes/v1"
	"github.com/josephmakin/monerohub/services"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error loading .env file: %v", err)
    }

    err = services.InitMonero(os.Getenv("RPC_HOST"), os.Getenv("RPC_PORT"))
    if err != nil {
        log.Fatalf("Error connecting to monero-wallet-rpc: %v", err)

    }

	collections := map[string]string{
		"payments": "payments",
	}
    err = services.InitMongo(os.Getenv("MONGO_URI"), "monerohub", collections)
    if err != nil {
        log.Fatalf("Error connecting to mongo: %v", err)
    }

    err = services.InitRedis(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"))
    if err != nil {
        log.Fatalf("Error connecting to redis: %v", err)
    }

    router := gin.Default()
    v1.SetupRoutes(router)

    err = router.Run(":8080")
    if err != nil {
        panic(err)
    }
}
