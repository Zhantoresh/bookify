package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bookify/internal/database"
	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/handlers"
	"github.com/bookify/internal/middleware"
	"github.com/bookify/internal/notification"
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
	userRepo := repository.NewUserRepository(db)
	notifier := notification.NewAsyncNotifier(notification.NewLogSender(log.Default()), 2, 32, log.Default())
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := notifier.Close(shutdownCtx); err != nil {
			log.Printf("Failed to shutdown notifier: %v", err)
		}
	}()

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize services
	specialistService := service.NewSpecialistService(specialistRepo, timeSlotRepo)
	bookingService := service.NewBookingService(bookingRepo, timeSlotRepo, specialistRepo, userRepo, notifier)
	timeSlotService := service.NewTimeSlotService(timeSlotRepo, userRepo, notifier)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userUsecase)
	bookingHandler := handlers.NewHandler(specialistService, bookingService)
	timeSlotHandler := handlers.NewTimeSlotHandler(timeSlotService)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Public routes (no auth required)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/specialists", bookingHandler.GetSpecialists)
	mux.HandleFunc("/specialistsWithSlots/", bookingHandler.GetSpecialistByID)

	// Protected routes (require auth)
	mux.Handle("/bookings", middleware.AuthMiddleware(http.HandlerFunc(bookingHandler.HandleBookings)))
	mux.Handle("/bookings/", middleware.AuthMiddleware(http.HandlerFunc(bookingHandler.HandleBookingsByID)))

	// Time slot management routes (specialist only - require auth and role check)
	specialistOnly := middleware.RoleMiddleware(string(domain.RoleSpecialist))
	mux.Handle("/time-slots", middleware.AuthMiddleware(specialistOnly(http.HandlerFunc(timeSlotHandler.HandleTimeSlots))))
	mux.Handle("/time-slots/", middleware.AuthMiddleware(specialistOnly(http.HandlerFunc(timeSlotHandler.HandleTimeSlotsWithID))))

	// Admin routes (require auth and admin role check)
	adminOnly := middleware.RoleMiddleware(string(domain.RoleAdmin))
	mux.Handle("/admin/dashboard", middleware.AuthMiddleware(adminOnly(http.HandlerFunc(handlers.AdminDashboard))))

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
