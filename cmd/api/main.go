package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bookify/internal/config"
	"github.com/bookify/internal/repository/postgres"
	appservice "github.com/bookify/internal/service"
	authsvc "github.com/bookify/internal/service/auth"
	httptransport "github.com/bookify/internal/transport/http"
	"github.com/bookify/internal/worker"
	"github.com/bookify/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	db, err := postgres.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Error("database_connection_failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	serviceRepo := postgres.NewServiceRepository(db)
	appointmentRepo := postgres.NewAppointmentRepository(db)

	jwtService := authsvc.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration)
	authService := appservice.NewAuthService(userRepo, jwtService)
	userService := appservice.NewUserService(userRepo)
	serviceService := appservice.NewServiceService(serviceRepo, userRepo)
	appointmentService := appservice.NewAppointmentService(appointmentRepo, serviceRepo, userRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	reminderWorker := worker.NewReminderWorker(appointmentRepo, log)
	reminderWorker.Start(ctx, &wg)

	workerPool := worker.NewWorkerPool(5, 100, log)
	workerPool.Start(ctx, &wg)

	handler := httptransport.NewServer(authService, userService, serviceService, appointmentService, jwtService, workerPool, log)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Info("server_starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server_failed", "error", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("server_shutting_down")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server_shutdown_failed", "error", err)
	}

	workerPool.Shutdown()
	wg.Wait()
	log.Info("server_stopped")
}
