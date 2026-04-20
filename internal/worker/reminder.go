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
}

func NewReminderWorker(repo repository.AppointmentRepository, logger *slog.Logger) *ReminderWorker {
	return &ReminderWorker{
		appointmentRepo: repo,
		logger:          logger,
		interval:        time.Hour,
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
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	appointments, err := w.appointmentRepo.GetAppointmentsByDateRange(ctx, start, end)
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
