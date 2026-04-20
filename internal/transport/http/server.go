package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/service"
	authsvc "github.com/bookify/internal/service/auth"
	"github.com/bookify/internal/transport/http/handler"
	"github.com/bookify/internal/transport/http/middleware"
	"github.com/bookify/internal/worker"
)

func NewServer(
	authService service.AuthService,
	userService service.UserService,
	serviceService service.ServiceService,
	appointmentService service.AppointmentService,
	jwtService *authsvc.JWTService,
	workerPool *worker.WorkerPool,
	logger *slog.Logger,
) nethttp.Handler {
	mux := nethttp.NewServeMux()

	authHandler := handler.NewAuthHandler(authService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService, workerPool, logger)
	userHandler := handler.NewUserHandler(userService)

	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)

	authMiddleware := middleware.AuthMiddleware(jwtService)
	protected := func(h nethttp.Handler) nethttp.Handler { return authMiddleware(h) }
	withRole := func(h nethttp.Handler, roles ...string) nethttp.Handler {
		return protected(middleware.RequireRole(roles...)(h))
	}

	mux.Handle("/api/v1/services", nethttp.HandlerFunc(serviceHandler.HandleCollection(withRole)))
	mux.Handle("/api/v1/services/my", withRole(nethttp.HandlerFunc(serviceHandler.ListMine), string(domain.RoleProvider)))
	mux.HandleFunc("/api/v1/services/", serviceHandler.HandleByID(protected, withRole))

	mux.Handle("/api/v1/auth/validate", protected(nethttp.HandlerFunc(authHandler.Validate)))
	mux.Handle("/api/v1/users/me", protected(nethttp.HandlerFunc(userHandler.Me)))

	mux.Handle("/api/v1/appointments", appointmentHandler.HandleCollection(protected, withRole))
	mux.Handle("/api/v1/appointments/my", protected(nethttp.HandlerFunc(appointmentHandler.ListMine)))
	mux.HandleFunc("/api/v1/appointments/available-slots", appointmentHandler.AvailableSlots)
	mux.Handle("/api/v1/appointments/", appointmentHandler.HandleByID(protected))

	return chain(mux,
		middleware.Recovery(logger),
		middleware.Logging(logger),
		middleware.CORS(),
	)
}

func chain(next nethttp.Handler, middlewares ...func(nethttp.Handler) nethttp.Handler) nethttp.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}
