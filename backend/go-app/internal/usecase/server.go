package usecase

import (
	"autofort/internal/entity"

	"github.com/google/uuid"
)

type Postgres interface {
	AddVehicleType(*entity.VehicleType) error

	AddVehicle(*entity.Vehicle) error
	GetVehicle(id uuid.UUID) (*entity.Vehicle, error)
	UpdateVehicle(*entity.Vehicle) error

	AddService(*entity.Service) error
	GetService(id uuid.UUID) (*entity.Service, error)
	SetServiceTime(time float64, serviceId, vehicleId uuid.UUID) error
}

type Server struct {
	CustomerRepo    CustomerRepo
	VehicleTypeRepo VehicleTypeRepo
	VehicleRepo     VehicleRepo
	WorkOrderClient WorkOrderClient
}

func NewService(CustomerRepo CustomerRepo, VehicleTypeRepo VehicleTypeRepo, VehicleRepo VehicleRepo, WorkOrderClient WorkOrderClient) *Server {
	return &Server{CustomerRepo: CustomerRepo, VehicleTypeRepo: VehicleTypeRepo, VehicleRepo: VehicleRepo, WorkOrderClient: WorkOrderClient}
}
