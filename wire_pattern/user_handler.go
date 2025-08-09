package main

import (
	"fmt"
	"net/http"
)

// HTTP handlers
type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response (simplified)
	fmt.Fprintf(w, `{"id": "%s", "name": "%s"}`, user.ID, user.Name)
}
