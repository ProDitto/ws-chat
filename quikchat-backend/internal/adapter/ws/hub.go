package ws

import (
	"log/slog"
	"quikchat/internal/domain"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	// Maps userID to a set of clients (for multiple connections per user)
	userClients map[int64]map[*Client]bool
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[int64]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if h.userClients[client.UserID] == nil {
				h.userClients[client.UserID] = make(map[*Client]bool)
			}
			h.userClients[client.UserID][client] = true
			h.mu.Unlock()
			slog.Info("client registered", "userID", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				if userClients, userExists := h.userClients[client.UserID]; userExists {
					delete(userClients, client)
					if len(userClients) == 0 {
						delete(h.userClients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			slog.Info("client unregistered", "userID", client.UserID)
		}
	}
}

// Register returns the channel for registering new clients.
func (h *Hub) Register() chan<- *Client {
	return h.register
}

// Unregister returns the channel for unregistering clients.
func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}

// BroadcastToUsers sends a message to all connected clients for the given user IDs.
func (h *Hub) BroadcastToUsers(message *domain.Message, userIDs []int64) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userID := range userIDs {
		if clients, ok := h.userClients[userID]; ok {
			for client := range clients {
				client.send(message)
			}
		}
	}
}

