package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type ReminderWorker struct {
	appointmentRepo repository.AppointmentRepository
	logger          *slog.Logger
	interval        time.Duration
	location        *time.Location
}

func NewReminderWorker(repo repository.AppointmentRepository, logger *slog.Logger, location *time.Location) *ReminderWorker {
	if location == nil {
		location = time.UTC
	}
	return &ReminderWorker{
		appointmentRepo: repo,
		logger:          logger,
		interval:        time.Hour,
		location:        location,
	}
}

func (w *ReminderWorker) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		w.logger.Info("reminder_worker_started")
		for {
			select {
			case <-ctx.Done():
				w.logger.Info("reminder_worker_stopping")
				return
			case <-ticker.C:
				w.checkAndSendReminders(ctx)
			}
		}
	}()
}

func (w *ReminderWorker) checkAndSendReminders(ctx context.Context) {
	tomorrow := time.Now().In(w.location).AddDate(0, 0, 1)
	startLocal := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, w.location)
	endLocal := startLocal.Add(24 * time.Hour)
	appointments, err := w.appointmentRepo.GetAppointmentsByDateRange(ctx, startLocal.UTC(), endLocal.UTC())
	if err != nil {
		w.logger.Error("reminder_fetch_failed", "error", err)
		return
	}
	for _, appointment := range appointments {
		if appointment.Status != domain.AppointmentConfirmed {
			continue
		}
		w.logger.Info("reminder_sent",
			"client_email", appointment.ClientEmail,
			"appointment_id", appointment.ID,
			"start_time", appointment.StartTime.Format(time.RFC3339),
		)
	}
}
