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

	"github.com/gorilla/mux"
	"github.com/srgjo27/e-learning/internal/entity"
	"github.com/srgjo27/e-learning/internal/infrastructure/repository"
	"github.com/srgjo27/e-learning/internal/interface/rest"
	"github.com/srgjo27/e-learning/internal/usecase"
	"github.com/srgjo27/e-learning/internal/utils"
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
	profileHandler := rest.NewProfileHandler(authUseCase)

	router := mux.NewRouter()

	router.HandleFunc("/v1/auth/register", authHandler.HandleRegister)
	router.HandleFunc("/v1/auth/login", authHandler.HandleLogin)
	router.HandleFunc("/v1/auth/password-reset/request", authHandler.HandlePasswordResetRequest)
	router.HandleFunc("/v1/auth/password-reset/reset", authHandler.HandlePasswordReset)

	router.Handle("v1/profile", utils.JWTMiddleware(authUseCase, profileHandler))

	adminHandler := rest.NewAdminHandler(authUseCase)

	adminSubrouter := router.PathPrefix("/v1/admin").Subrouter()
	adminSubrouter.Use(func(next http.Handler) http.Handler {
		return utils.JWTMiddleware(authUseCase, utils.RBACMiddleware(entity.RoleAdmin)(next))
	})
	adminSubrouter.HandleFunc("/users", adminHandler.ListUsers).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/users/{id}/role", adminHandler.UpdateUserRole).Methods(http.MethodPut)
	adminSubrouter.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods(http.MethodDelete)


	router.Handle("/teacher-area", utils.JWTMiddleware(authUseCase, utils.RBACMiddleware(entity.RoleTeacher)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Teacher Area"))
	}))))
	router.Handle("/student-area", utils.JWTMiddleware(authUseCase, utils.RBACMiddleware(entity.RoleStudent)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Student Area"))
	}))))

	srv := &http.Server{
		Addr:	":8080",
		Handler: router,
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