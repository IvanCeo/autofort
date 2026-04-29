package entity

import (
	appval "autofort/internal/validator"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID
	FirstName   string `validate:"required"`
	LastName    string `validate:"required"`
	PhoneNumber string `validate:"required,e164"`
}

func (c *Customer) Validate() error {
	return appval.Validate.Struct(c)
}
