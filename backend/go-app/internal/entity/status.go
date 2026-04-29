package entity

import (
	"autofort/internal/validator"

	"github.com/google/uuid"
)

type Status struct {
	ID   uuid.UUID `validate:"required,uuid"`
	Name string    `validate:"required"`
}

func (s *Status) Validate() error {
	return validator.Validate.Struct(s)
}
