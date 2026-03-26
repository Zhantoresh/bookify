package main

import (
	"log"
	"net/http"
	"os"	
	
	"github.com/bookify/internal/database"
	"github.com/bookify/internal/handlers"
	"github.com/bookify/internal/repository"
	"github.com/bookify/internal/service"
  "github.com/bookify/internal/usecase"
)

func main() {
	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "bookify"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database
	db, err := database.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	specialistRepo := repository.NewSpecialistRepository(db)
	timeSlotRepo := repository.NewTimeSlotRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
  
 	userRepo := database.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	authHandler := handlers.NewAuthHandler(userUsecase)

	// Initialize services
	specialistService := service.NewSpecialistService(specialistRepo, timeSlotRepo)
	bookingService := service.NewBookingService(bookingRepo, timeSlotRepo, specialistRepo)

	// Initialize handlers
	handler := handlers.NewHandler(specialistService, bookingService)

	// Setup HTTP routes
	mux := http.NewServeMux()
  
  mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
  
	handler.RegisterRoutes(mux)

	// Start server
	addr := getEnv("SERVER_ADDR", ":8080")
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
