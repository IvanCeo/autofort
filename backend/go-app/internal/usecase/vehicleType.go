package usecase

import (
	"errors"
	"strings"

	"autofort/internal/entity"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type VehicleTypeRepo interface {
	ListVehicleTypes() ([]*entity.VehicleType, error)
	AddVehicleType(*entity.VehicleType) error
	GetVehicleTypeByID(uuid.UUID) (*entity.VehicleType, error)
}

var (
	ErrVehicleTypeBrandRequired = errors.New("vehicle type brand is required")
	ErrVehicleTypeModelRequired = errors.New("vehicle type model is required")
)

var brandCaser = cases.Title(language.Und)

func (s *Server) AddVehicleType(brand, model string) (uuid.UUID, error) {
	brand = strings.TrimSpace(brand)
	model = strings.TrimSpace(model)

	if brand == "" {
		return uuid.Nil, ErrVehicleTypeBrandRequired
	}
	if model == "" {
		return uuid.Nil, ErrVehicleTypeModelRequired
	}

	brand = brandCaser.String(strings.ToLower(brand))

	v := &entity.VehicleType{
		ID:    uuid.New(),
		Brand: brand,
		Model: model,
	}

	if err := v.Validate(); err != nil {
		return uuid.Nil, err
	}

	if err := s.VehicleTypeRepo.AddVehicleType(v); err != nil {
		return uuid.Nil, err
	}

	return v.ID, nil
}

func (s *Server) ListVehicleTypes() ([]*entity.VehicleType, error) {
	return s.VehicleTypeRepo.ListVehicleTypes()
}

func (s *Server) GetVehicleTypeByID(id uuid.UUID) (*entity.VehicleType, error) {
	return s.VehicleTypeRepo.GetVehicleTypeByID(id)
}
