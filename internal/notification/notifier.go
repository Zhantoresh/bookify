package notification

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Message struct {
	Recipient string
	Subject   string
	Body      string
}

type Sender interface {
	Send(ctx context.Context, message Message) error
}

type Notifier interface {
	Notify(message Message)
	Close(ctx context.Context) error
}

type AsyncNotifier struct {
	sender Sender
	logger *log.Logger
	queue  chan Message
	wg     sync.WaitGroup
	once   sync.Once
}

func NewAsyncNotifier(sender Sender, workers, bufferSize int, logger *log.Logger) *AsyncNotifier {
	if workers < 1 {
		workers = 1
	}
	if bufferSize < 1 {
		bufferSize = 1
	}
	if logger == nil {
		logger = log.Default()
	}

	n := &AsyncNotifier{
		sender: sender,
		logger: logger,
		queue:  make(chan Message, bufferSize),
	}

	for i := 0; i < workers; i++ {
		n.wg.Add(1)
		go n.worker()
	}

	return n
}

func (n *AsyncNotifier) Notify(message Message) {
	select {
	case n.queue <- message:
	default:
		n.logger.Printf("notification queue is full, dropping message for %s", message.Recipient)
	}
}

func (n *AsyncNotifier) Close(ctx context.Context) error {
	n.once.Do(func() {
		close(n.queue)
	})

	done := make(chan struct{})
	go func() {
		n.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (n *AsyncNotifier) worker() {
	defer n.wg.Done()

	for message := range n.queue {
		sendCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := n.sender.Send(sendCtx, message); err != nil {
			n.logger.Printf("failed to send notification to %s: %v", message.Recipient, err)
		}
		cancel()
	}
}

type LogSender struct {
	logger *log.Logger
}

func NewLogSender(logger *log.Logger) *LogSender {
	if logger == nil {
		logger = log.Default()
	}

	return &LogSender{logger: logger}
}

func (s *LogSender) Send(_ context.Context, message Message) error {
	s.logger.Printf("notification sent to=%s subject=%q body=%q", message.Recipient, message.Subject, message.Body)
	return nil
}

type NoopNotifier struct{}

func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{}
}

func (n *NoopNotifier) Notify(message Message) {}

func (n *NoopNotifier) Close(ctx context.Context) error {
	return nil
}

func BuildTimeSlotCreatedMessage(recipient, name string, slotTime time.Time) Message {
	return Message{
		Recipient: recipient,
		Subject:   "New time slot created",
		Body:      fmt.Sprintf("Hello, %s. Your new time slot for %s has been created.", name, slotTime.Format(time.RFC3339)),
	}
}

func BuildBookingCreatedMessage(recipient, name, specialist string, slotTime time.Time) Message {
	return Message{
		Recipient: recipient,
		Subject:   "Booking confirmed",
		Body:      fmt.Sprintf("Hello, %s. Your booking with %s for %s is confirmed.", name, specialist, slotTime.Format(time.RFC3339)),
	}
}
