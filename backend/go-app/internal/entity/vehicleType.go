package entity

import (
	"autofort/internal/validator"

	"github.com/google/uuid"
)

type VehicleType struct {
	ID    uuid.UUID `validate:"required"`
	Brand string    `validate:"required"`
	Model string    `validate:"required"`
}

func (v *VehicleType) Validate() error {
	return validator.Validate.Struct(v)
}
