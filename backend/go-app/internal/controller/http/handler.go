package http

import "autofort/internal/usecase"

type Handler struct {
	server *usecase.Server
}

func NewHandler(s *usecase.Server) *Handler {
	return &Handler{server: s}
}
