package main

import (
	"context"
	"log"
	"log/slog"
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


	// Logger initialization
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, 
	}))
	slog.SetDefault(logger)

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
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Database connection established")

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
			logger.Error("Failed to shutdown notifier", "error", err)
		}
	}()

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, logger)

	// Initialize services
	specialistService := service.NewSpecialistService(specialistRepo, timeSlotRepo)
	bookingService := service.NewBookingService(bookingRepo, timeSlotRepo, specialistRepo, userRepo, notifier, logger)
	timeSlotService := service.NewTimeSlotService(timeSlotRepo, userRepo, notifier, logger)

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


	finalHandler := middleware.LoggingMiddleware(logger)(mux)
	
	// Start server
	addr := getEnv("SERVER_ADDR", ":8080")
	logger.Info("Starting server", "addr", addr)

	if err := http.ListenAndServe(addr, finalHandler); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
