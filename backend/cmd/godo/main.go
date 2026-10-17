// Command godo runs the GoDo API server.
//
// Configuration comes from environment variables (or a .env file); see
// internal/config and .env.example.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/NathanGrs00/godo/internal/api"
	"github.com/NathanGrs00/godo/internal/auth"
	"github.com/NathanGrs00/godo/internal/config"
	"github.com/NathanGrs00/godo/internal/realtime"
	mongorepo "github.com/NathanGrs00/godo/internal/repository/mongo"
	"github.com/NathanGrs00/godo/internal/service"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Stop gracefully on Ctrl+C or when Docker stops the container
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := connectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	log.Println("Connected to MongoDB")

	db := client.Database(cfg.MongoDatabase)
	if err := mongorepo.EnsureIndexes(ctx, db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           newRouter(cfg, db),
		ReadHeaderTimeout: 10 * time.Second,
		// No WriteTimeout: the /events stream stays open for as long as the app is open
	}

	go func() {
		log.Printf("Listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}
}

// connectMongo opens and pings a MongoDB connection.
func connectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(connectCtx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

// newRouter wires repositories, services and handlers together.
func newRouter(cfg config.Config, db *mongo.Database) http.Handler {
	hub := realtime.NewHub()
	tokens := auth.NewTokenManager(cfg.JWTSecret, cfg.TokenTTL)

	taskRepo := mongorepo.NewTaskRepository(db)
	deadlineRepo := mongorepo.NewDeadlineRepository(db)
	tagRepo := mongorepo.NewTagRepository(db)
	userRepo := mongorepo.NewUserRepository(db)

	tasks := service.NewTaskService(taskRepo, deadlineRepo, tagRepo, hub)
	deadlines := service.NewDeadlineService(deadlineRepo, tagRepo, tasks, hub)
	tags := service.NewTagService(tagRepo, tasks, deadlines, hub)

	return api.NewRouter(api.Deps{
		Auth:           service.NewAuthService(userRepo, tokens),
		Tasks:          tasks,
		Deadlines:      deadlines,
		Tags:           tags,
		Tokens:         tokens,
		Events:         hub,
		AllowedOrigins: cfg.AllowedOrigins,
	})
}
