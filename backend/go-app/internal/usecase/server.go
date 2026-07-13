package usecase

import (
	"autofort/internal/entity"
	er "autofort/internal/errors"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/google/uuid"
)

type Postgres interface {
	AddVehicleType(ctx context.Context, e *entity.VehicleType) error

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
	UserRepo        UserRepo

	key *ecdsa.PrivateKey
}

func NewService(CustomerRepo CustomerRepo, VehicleTypeRepo VehicleTypeRepo, VehicleRepo VehicleRepo, WorkOrderClient WorkOrderClient, UserRepo UserRepo) (*Server, error) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, er.ErrGenerateKey
	}
	return &Server{
		CustomerRepo:    CustomerRepo,
		VehicleTypeRepo: VehicleTypeRepo,
		VehicleRepo:     VehicleRepo,
		WorkOrderClient: WorkOrderClient,
		UserRepo:        UserRepo,
		key:             k,
	}, nil
}
