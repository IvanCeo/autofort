package entity

import "github.com/google/uuid"

type Order struct {
	ID         uuid.UUID `validate:"required,uuid"`
	CustomerId uuid.UUID `validate:"required,uuid"`
	VehicleId  uuid.UUID `validate:"required,uuid"`
	PartId     uuid.UUID `validate:"required,uuid"`
	TotalCost  int       `validate:"gte=0"`
}
