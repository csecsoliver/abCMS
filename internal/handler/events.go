package handler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SSEMessage represents a named Server-Sent Event payload.
type SSEMessage struct {
	Event string
	Data  string
}

// EventBroker fans out Server-Sent Events to connected HTMX clients.
type EventBroker struct {
	mu      sync.RWMutex
	clients map[chan SSEMessage]struct{}
	closed  bool
}

// NewEventBroker creates an initialized SSE broker.
func NewEventBroker() *EventBroker {
	return &EventBroker{
		clients: make(map[chan SSEMessage]struct{}),
	}
}

// Subscribe registers a new client channel and returns an unsubscribe function.
func (b *EventBroker) Subscribe() (<-chan SSEMessage, func()) {
	ch := make(chan SSEMessage, 16)
	b.mu.Lock()
	if b.closed {
		close(ch)
		b.mu.Unlock()
		return ch, func() {}
	}
	b.clients[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	unsub := func() {
		once.Do(func() {
			b.mu.Lock()
			if _, ok := b.clients[ch]; ok {
				delete(b.clients, ch)
				close(ch)
			}
			b.mu.Unlock()
		})
	}
	return ch, unsub
}

// Publish sends an SSE event to all connected subscribers without blocking slow clients.
func (b *EventBroker) Publish(event, data string) {
	msg := SSEMessage{Event: event, Data: data}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return
	}
	for ch := range b.clients {
		select {
		case ch <- msg:
		default:
			// Drop if client buffer is full so mutations never block.
		}
	}
}

// Close disconnects all active SSE subscribers so graceful HTTP shutdown never hangs.
func (b *EventBroker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for ch := range b.clients {
		delete(b.clients, ch)
		close(ch)
	}
}

func (h *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch, unsub := h.events.Subscribe()
	defer unsub()

	// Send an initial comment so the browser opens the EventSource stream immediately.
	_, _ = fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	keepalive := time.NewTicker(25 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg, open := <-ch:
			if !open {
				return
			}
			writeSSEMessage(w, msg)
			flusher.Flush()
		case <-keepalive.C:
			_, _ = fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

func writeSSEMessage(w http.ResponseWriter, msg SSEMessage) {
	if msg.Event != "" {
		_, _ = fmt.Fprintf(w, "event: %s\n", msg.Event)
	}
	lines := strings.Split(strings.ReplaceAll(msg.Data, "\r\n", "\n"), "\n")
	for _, line := range lines {
		_, _ = fmt.Fprintf(w, "data: %s\n", line)
	}
	_, _ = fmt.Fprint(w, "\n")
}
