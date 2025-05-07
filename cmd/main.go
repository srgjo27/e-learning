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

	jwtSecret := "supersecretkey"
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
	courseCollection := client.Database("e-learning").Collection("courses")
	classCollection := client.Database("e-learning").Collection("classes")
	announcementCollection := client.Database("e-learning").Collection("announcements")
	assignmentCollection := client.Database("e-learning").Collection("assignments")
	assessmentCollection := client.Database("e-learning").Collection("assessments")
	messageCollection := client.Database("e-learning").Collection("messages")

	userRepo := repository.NewMongoUserRepository(userCollection)
	courseRepo := repository.NewMongoCourseRepository(courseCollection)
	classRepo := repository.NewMongoClassRepository(classCollection)
	announcementRepo := repository.NewMongoAnnouncementRepository(announcementCollection)
	assignmentRepo := repository.NewMongoAssignmentRepository(assignmentCollection)
	assessmentRepo := repository.NewMongoAssessmentRepository(assessmentCollection)
	messageRepo := repository.NewMongoMessageRepository(messageCollection)

	authUseCase := usecase.NewAuthUseCase(userRepo, []byte(jwtSecret))
	adminUseCase := usecase.NewAdminUseCase(courseRepo, classRepo, announcementRepo)
	teacherUseCase := usecase.NewTeacherUseCase(courseRepo, classRepo, userRepo)
	teacherAdvancedUseCase := usecase.NewTeacherAdvancedUseCase(assignmentRepo, assessmentRepo, messageRepo)

	authHandler := rest.NewAuthHandler(authUseCase)
	profileHandler := rest.NewProfileHandler(authUseCase)
	adminTasksHandler := rest.NewAdminTasksHandler(adminUseCase)
	adminHandler := rest.NewAdminHandler(authUseCase)
	teacherHandler := rest.NewTeacherHandler(teacherUseCase)
	teacherAdvancedHandler := rest.NewTeacherAdvancedHandler(teacherAdvancedUseCase)

	router := mux.NewRouter()

	router.HandleFunc("/v1/auth/register", authHandler.HandleRegister)
	router.HandleFunc("/v1/auth/login", authHandler.HandleLogin)
	router.HandleFunc("/v1/auth/password-reset/request", authHandler.HandlePasswordResetRequest)
	router.HandleFunc("/v1/auth/password-reset/reset", authHandler.HandlePasswordReset)

	router.Handle("v1/profile", utils.JWTMiddleware(authUseCase, profileHandler))

	adminSubrouter := router.PathPrefix("/v1/admin").Subrouter()
	adminSubrouter.Use(func(next http.Handler) http.Handler {
		return utils.JWTMiddleware(authUseCase, utils.RBACMiddleware(entity.RoleAdmin)(next))
	})
	adminSubrouter.HandleFunc("/users", adminHandler.ListUsers).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/users/{id}/role", adminHandler.UpdateUserRole).Methods(http.MethodPut)
	adminSubrouter.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods(http.MethodDelete)

	adminSubrouter.HandleFunc("/courses", adminTasksHandler.ListCourses).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/courses", adminTasksHandler.CreateCourse).Methods(http.MethodPost)
	adminSubrouter.HandleFunc("/courses/{id}", adminTasksHandler.GetCourse).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/courses/{id}", adminTasksHandler.UpdateCourse).Methods(http.MethodPut)
	adminSubrouter.HandleFunc("/courses/{id}", adminTasksHandler.DeleteCourse).Methods(http.MethodDelete)

	adminSubrouter.HandleFunc("/classes", adminTasksHandler.ListClasses).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/classes", adminTasksHandler.CreateClass).Methods(http.MethodPost)
	adminSubrouter.HandleFunc("/classes/{id}", adminTasksHandler.GetClass).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/classes/{id}", adminTasksHandler.UpdateClass).Methods(http.MethodPut)
	adminSubrouter.HandleFunc("/classes/{id}", adminTasksHandler.DeleteClass).Methods(http.MethodDelete)

	adminSubrouter.HandleFunc("/announcements", adminTasksHandler.ListAnnouncements).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/announcements", adminTasksHandler.CreateAnnouncement).Methods(http.MethodPost)
	adminSubrouter.HandleFunc("/announcements/{id}", adminTasksHandler.GetAnnouncement).Methods(http.MethodGet)
	adminSubrouter.HandleFunc("/announcements/{id}", adminTasksHandler.UpdateAnnouncement).Methods(http.MethodPut)
	adminSubrouter.HandleFunc("/announcements/{id}", adminTasksHandler.DeleteAnnouncement).Methods(http.MethodDelete)


	teacherSubrouter := router.PathPrefix("/v1/teacher").Subrouter()
	teacherSubrouter.Use(func(next http.Handler) http.Handler {
		return utils.JWTMiddleware(authUseCase, utils.RBACMiddleware(entity.RoleTeacher)(next))
	})
	teacherSubrouter.HandleFunc("/courses", teacherHandler.ListCourses).Methods(http.MethodGet)
	teacherSubrouter.HandleFunc("/classes", teacherHandler.ListClasses).Methods(http.MethodGet)
	teacherSubrouter.HandleFunc("/classes/{id}/students", teacherHandler.ListStudents).Methods(http.MethodGet)

	teacherSubrouter.HandleFunc("/assignments", teacherAdvancedHandler.ListAssignments).Methods(http.MethodGet)
	teacherSubrouter.HandleFunc("/assignments", teacherAdvancedHandler.CreateAssignment).Methods(http.MethodPost)
	// Optional: Implement GET, PUT, DELETE for /assignments/{id} similarly

	teacherSubrouter.HandleFunc("/messages", teacherAdvancedHandler.ListMessages).Methods(http.MethodGet)
	teacherSubrouter.HandleFunc("/messages", teacherAdvancedHandler.CreateMessage).Methods(http.MethodPost)
	
	// router.Handle("/student-area", utils.JWTMiddleware(authUseCase, utils.RBACMiddleware(entity.RoleStudent)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome to Student Area"))
	// }))))

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