package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"quikchat/internal/adapter/ws"
	"quikchat/internal/domain"
	"quikchat/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for now. In production, you'd want to restrict this.
		return true
	},
}

type MessageHandler struct {
	messageService usecase.MessageUseCase
	hub            *ws.Hub
}

func NewMessageHandler(messageService usecase.MessageUseCase, hub *ws.Hub) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		hub:            hub,
	}
}

type SendMessageRequest struct {
	Type    domain.MessageType `json:"type"`
	Content string             `json:"content"`
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationIDStr := chi.URLParam(r, "conversationID")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := h.messageService.SendMessage(r.Context(), userID, conversationID, req.Type, req.Content)
	if err != nil {
		slog.Error("failed to send message", "error", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationIDStr := chi.URLParam(r, "conversationID")
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	cursorStr := r.URL.Query().Get("cursor")
	cursor, _ := strconv.ParseInt(cursorStr, 10, 64) // Defaults to 0 if empty or invalid

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr) // Defaults to 0 if empty or invalid

	messages, err := h.messageService.GetMessagesByConversationID(r.Context(), userID, conversationID, cursor, limit)
	if err != nil {
		slog.Error("failed to get messages", "error", err)
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *MessageHandler) GetPresignedURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename query parameter is required", http.StatusBadRequest)
		return
	}

	uploadURL, key, err := h.messageService.GetPresignedUploadURL(r.Context(), userID, filename)
	if err != nil {
		slog.Error("failed to get presigned url", "error", err)
		http.Error(w, "Could not generate upload URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"upload_url": uploadURL,
		"key":        key,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MessageHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		// Cannot upgrade connection if user is not authenticated
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade websocket", "error", err)
		return
	}

	client := &ws.Client{
		UserID: userID,
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	client.Hub.Register() <- client

	go client.WritePump()
	go client.ReadPump()
}

