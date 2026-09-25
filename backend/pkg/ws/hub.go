package ws

import (
	"encoding/json"
	"log/slog"
	"sync"
)

// Message 下发给客户端的 WebSocket 消息。
type Message struct {
	Type       string `json:"type"`
	ActivityID int64  `json:"activity_id"`
	Data       any    `json:"data,omitempty"`
}

// Hub 维护活动排行榜实时推送的连接。
type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[chan []byte]struct{}
	logger  *slog.Logger
}

// NewHub 构造 Hub。
func NewHub(logger *slog.Logger) *Hub {
	return &Hub{clients: make(map[int64]map[chan []byte]struct{}), logger: logger}
}

// Subscribe 订阅某活动的更新通道。
func (h *Hub) Subscribe(activityID int64) chan []byte {
	ch := make(chan []byte, 16)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[activityID] == nil {
		h.clients[activityID] = make(map[chan []byte]struct{})
	}
	h.clients[activityID][ch] = struct{}{}
	return ch
}

// Unsubscribe 取消订阅。
func (h *Hub) Unsubscribe(activityID int64, ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[activityID]; ok {
		delete(set, ch)
		close(ch)
		if len(set) == 0 {
			delete(h.clients, activityID)
		}
	}
}

// Broadcast 广播消息到某活动的所有订阅者。
func (h *Hub) Broadcast(activityID int64, msg Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("ws marshal message", "err", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.clients[activityID] {
		select {
		case ch <- payload:
		default:
			h.logger.Warn("ws channel full, skip", "activity_id", activityID)
		}
	}
}
