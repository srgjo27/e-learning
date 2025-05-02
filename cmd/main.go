package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/srgjo27/e-learning/internal/infrastructure/repository"
	"github.com/srgjo27/e-learning/internal/interface/rest"
	"github.com/srgjo27/e-learning/internal/usecase"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	if uri := os.Getenv("DB_HOST"); uri != "" {
		mongoURI = uri
	}

	jwtSecret := ""
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = secret
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}
	log.Println("Connected to MongoDB")

	userCollection := client.Database("e-learning").Collection("users")

	userRepo := repository.NewMongoUserRepository(userCollection)

	authUseCase := usecase.NewAuthUseCase(userRepo, []byte(jwtSecret))

	authHandler := rest.NewAuthHandler(authUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/auth/register", authHandler.HandleRegister)
	mux.HandleFunc("/v1/auth/login", authHandler.HandleLogin)
	mux.HandleFunc("/v1/auth/password-reset/request", authHandler.HandlePasswordResetRequest)
	mux.HandleFunc("/v1/auth/password-reset/reset", authHandler.HandlePasswordReset)

	srv := &http.Server{
		Addr:	":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("MongoDB disconnet error: %v", err)
	}
	log.Println("MongoDB connection closed")
	log.Println("Server exited properly")
}