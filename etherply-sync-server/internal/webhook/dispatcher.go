package webhook

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// EventPayload represents the JSON body sent to the webhook URL.
type EventPayload struct {
	Event     string      `json:"event"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type Dispatcher struct {
	url    string
	queue  chan EventPayload
	client *http.Client
	logger *slog.Logger
}

func NewDispatcher(url string) *Dispatcher {
	// If no URL configured, we still return a Dispatcher but it basically no-ops safely.
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	d := &Dispatcher{
		url:    url,
		queue:  make(chan EventPayload, 1000), // Buffer events
		client: &http.Client{Timeout: 5 * time.Second},
		logger: logger,
	}

	if url != "" {
		go d.worker() // Start background worker
	}

	return d
}

func (d *Dispatcher) Dispatch(event string, payload interface{}) {
	if d.url == "" {
		return
	}

	// Non-blocking send. If queue full, drop event to prioritize server stability.
	select {
	case d.queue <- EventPayload{
		Event:     event,
		Timestamp: time.Now().UnixMicro(),
		Payload:   payload,
	}:
	default:
		d.logger.Warn("webhook_queue_full", slog.String("event", event))
	}
}

func (d *Dispatcher) worker() {
	for evt := range d.queue {
		d.send(evt)
	}
}

func (d *Dispatcher) send(evt EventPayload) {
	body, err := json.Marshal(evt)
	if err != nil {
		d.logger.Error("webhook_marshal_failed", slog.Any("error", err))
		return
	}

	resp, err := d.client.Post(d.url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		d.logger.Error("webhook_send_failed", slog.Any("error", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		d.logger.Warn("webhook_returned_error", slog.Int("status", resp.StatusCode))
	}
}
