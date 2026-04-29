package entity

import (
	"autofort/internal/validator"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID          uuid.UUID `validate:"required,uuid"`
	CustomerId  uuid.UUID `validate:"required,uuid"`
	VehicleType uuid.UUID `validate:"required,uuid"`
	Vin         string
	GovNumber   string `validate:"required"`
	Mileage     int    `validate:"required,gte=0"`
}

func (v *Vehicle) Validate() error {
	return validator.Validate.Struct(v)
}
