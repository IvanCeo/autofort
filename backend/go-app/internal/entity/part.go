package entity

import (
	"autofort/internal/validator"

	"github.com/google/uuid"
)

type Part struct {
	ID   uuid.UUID `validate:"required,uuid"`
	Name string
}

func (p *Part) Validate() error {
	return validator.Validate.Struct(p)
}
