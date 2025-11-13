package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/NathanGrs00/godo/internal/api"
	mongorepo "github.com/NathanGrs00/godo/internal/repository/mongo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, relying on system env variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGO_URI is empty — environment variable not loaded")
	}

	log.Println("✅ MONGO_URI loaded:", mongoURI[:10], "...")

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	db := client.Database("godo")

	taskRepo := mongorepo.NewTaskRepository(db)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api.SetupRoutes(r, taskRepo)

	r.Run(":8080")
}
