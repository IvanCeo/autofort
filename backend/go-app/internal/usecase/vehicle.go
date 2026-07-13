package usecase

import (
	"context"
	"errors"
	"strings"

	"autofort/internal/entity"
	er "autofort/internal/errors"

	"github.com/google/uuid"
)

var (
	ErrVehicleGovNumberRequired = errors.New("vehicle gov number is required")
	ErrVehicleTypeRequired      = errors.New("vehicle type is required")
)

type VehicleRepo interface {
	AddVehicle(*entity.Vehicle) error
	ListVehiclesByCustomerIDs(ids []uuid.UUID) ([]*entity.Vehicle, error)

	GetVehicle(id uuid.UUID) (*entity.Vehicle, error)
	UpdateVehicle(*entity.Vehicle) error
}

type UpdateVehicleInput struct {
	ID        uuid.UUID
	Vin       *string
	GovNumber *string
	Mileage   *int
}

func (s *Server) EditVehicle(update UpdateVehicleInput) error {
	v, err := s.VehicleRepo.GetVehicle(update.ID)
	if err != nil {
		return err
	}

	if update.Vin != nil {
		v.Vin = strings.TrimSpace(*update.Vin)
	}

	if update.GovNumber != nil {
		v.GovNumber = strings.TrimSpace(*update.GovNumber)
	}

	if update.Mileage != nil {
		v.Mileage = *update.Mileage
	}

	if err := v.Validate(); err != nil {
		return err
	}

	return s.VehicleRepo.UpdateVehicle(v)
}

func (s *Server) AddVehicleToCustomer(
	ctx context.Context,
	customerID uuid.UUID,
	vehicleTypeID uuid.UUID,
	vin string,
	govNumber string,
	mileage *int,
) (uuid.UUID, error) {

	if customerID == uuid.Nil {
		return uuid.Nil, er.ErrCustomerIDRequired
	}
	if vehicleTypeID == uuid.Nil {
		return uuid.Nil, ErrVehicleTypeRequired
	}

	govNumber = strings.TrimSpace(govNumber)
	if govNumber == "" {
		return uuid.Nil, ErrVehicleGovNumberRequired
	}

	vin = strings.TrimSpace(vin)

	finalMileage := 0
	if mileage != nil {
		finalMileage = *mileage
	}

	if _, err := s.CustomerRepo.GetCustomer(customerID); err != nil {
		return uuid.Nil, err
	}

	if _, err := s.VehicleTypeRepo.GetVehicleTypeByID(ctx, vehicleTypeID); err != nil {
		return uuid.Nil, err
	}

	v := &entity.Vehicle{
		ID:          uuid.New(),
		CustomerId:  customerID,
		VehicleType: vehicleTypeID,
		Vin:         vin,
		GovNumber:   govNumber,
		Mileage:     finalMileage,
	}

	if err := v.Validate(); err != nil {
		return uuid.Nil, err
	}

	if err := s.VehicleRepo.AddVehicle(v); err != nil {
		return uuid.Nil, err
	}

	return v.ID, nil
}
