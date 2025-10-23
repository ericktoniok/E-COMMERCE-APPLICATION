package realtime

import (
	"sync"
)

type Subscriber chan []byte

type Hub struct {
	mu   sync.RWMutex
	subs map[Subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{subs: make(map[Subscriber]struct{})}
}

func (h *Hub) Subscribe() Subscriber {
	ch := make(Subscriber, 8)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch Subscriber) {
	h.mu.Lock()
	if _, ok := h.subs[ch]; ok {
		delete(h.subs, ch)
		close(ch)
	}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(payload []byte) {
	h.mu.RLock()
	for ch := range h.subs {
		select {
		case ch <- payload:
		default:
			// drop if slow
		}
	}
	h.mu.RUnlock()
}
