package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"quikchat/internal/domain"
	"quikchat/internal/service"
	"quikchat/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	groupService usecase.GroupUseCase
}

func NewGroupHandler(groupService usecase.GroupUseCase) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

type CreateGroupRequest struct {
	Name        string  `json:"name"`
	Tag         string  `json:"tag"`
	Description *string `json:"description"`
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group, err := h.groupService.CreateGroup(r.Context(), userID, req.Name, req.Tag, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) SearchGroups(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	if tag == "" {
		http.Error(w, "tag query parameter is required", http.StatusBadRequest)
		return
	}

	groups, err := h.groupService.SearchGroups(r.Context(), tag, 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

func (h *GroupHandler) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.ParseInt(chi.URLParam(r, "groupID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	group, members, err := h.groupService.GetGroupDetails(r.Context(), userID, groupID)
	if err != nil {
		if errors.Is(err, service.ErrGroupNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		Group   *domain.Group         `json:"group"`
		Members []*domain.GroupMember `json:"members"`
	}{
		Group:   group,
		Members: members,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.ParseInt(chi.URLParam(r, "groupID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	err = h.groupService.JoinGroup(r.Context(), userID, groupID)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyInGroup) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.ParseInt(chi.URLParam(r, "groupID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	err = h.groupService.LeaveGroup(r.Context(), userID, groupID)
	if err != nil {
		if errors.Is(err, service.ErrOwnerCannotLeave) || errors.Is(err, service.ErrNotInGroup) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

