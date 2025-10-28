package handler

import (
	"encoding/json"
	"net/http"
	"quikchat/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FriendHandler struct {
	friendService usecase.FriendUseCase
}

func NewFriendHandler(friendService usecase.FriendUseCase) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

type SendFriendRequest struct {
	Username string `json:"username"`
}

type RespondToRequest struct {
	Action string `json:"action"` // "accept" or "decline"
}

func (h *FriendHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value("user_id").(int64)
	var req SendFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	friendReq, err := h.friendService.SendRequest(r.Context(), senderID, req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(friendReq)
}

func (h *FriendHandler) ListIncomingRequests(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	requests, err := h.friendService.ListIncomingRequests(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(requests)
}

func (h *FriendHandler) RespondToRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	requestID, err := strconv.ParseInt(chi.URLParam(r, "requestID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	var req RespondToRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.friendService.RespondToRequest(r.Context(), userID, requestID, req.Action); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FriendHandler) ListFriends(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	friends, err := h.friendService.ListFriends(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(friends)
}

func (h *FriendHandler) Unfriend(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	friendUsername := chi.URLParam(r, "username")

	if err := h.friendService.Unfriend(r.Context(), userID, friendUsername); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

