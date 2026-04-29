package entity

import (
	"autofort/internal/validator"

	"github.com/google/uuid"
)

type Service struct {
	ID   uuid.UUID `validate:"required,uuid"`
	Name string    `validate:"required"`
}

func (s *Service) Validate() error {
	return validator.Validate.Struct(s)
}
