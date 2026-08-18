package http

import (
	"autofort/internal/usecase"
)

type Handler struct {
	server *usecase.Server
}

func NewHandler(s *usecase.Server) (*Handler, error) {
	return &Handler{server: s}, nil
}

func (h *Handler) GetServer() *usecase.Server {
	return h.server
}

func (h *Handler) SignUp() error
