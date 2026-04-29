//go:build ignore
// +build ignore

package usecase

import (
	"autofort/internal/entity"

	"github.com/google/uuid"
)

func (s *Server) AddService(name string) error {
	service := &entity.Service{
		ID:   uuid.New(),
		Name: name,
	}
	if err := service.Validate(); err != nil {
		return err
	}

	return s.Postgres.AddService(service)
}

// func (s *Server) SetServiceTimeToVehicleType(vehicleType, serviceId uuid.UUID, time float64) error {

// }
