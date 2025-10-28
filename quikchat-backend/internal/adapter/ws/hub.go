package ws

import (
	"encoding/json"
	"log/slog"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Map of user IDs to their connected clients.
	userClients map[int64]map[*Client]bool

	// Mutex for protecting userClients map during broadcasts
	mu sync.RWMutex
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
			slog.Info("Client registered", "userID", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				if h.userClients[client.UserID] != nil {
					delete(h.userClients[client.UserID], client)
					if len(h.userClients[client.UserID]) == 0 {
						delete(h.userClients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			slog.Info("Client unregistered", "userID", client.UserID)
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

// BroadcastToUsers sends a message to all connected clients for a given list of user IDs.
func (h *Hub) BroadcastToUsers(payload interface{}, userIDs []int64) {
	message, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal broadcast message", "error", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userID := range userIDs {
		if clients, ok := h.userClients[userID]; ok {
			for client := range clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
