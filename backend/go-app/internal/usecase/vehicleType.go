package usecase

import (
	"context"
	"errors"
	"strings"

	"autofort/internal/entity"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type VehicleTypeRepo interface {
	ListVehicleTypes(ctx context.Context) ([]*entity.VehicleType, error)
	GetVehicleTypeByID(ctx context.Context, u uuid.UUID) (*entity.VehicleType, error)
	AddVehicleType(ctx context.Context, e *entity.VehicleType) error
}

var (
	ErrVehicleTypeBrandRequired = errors.New("vehicle type brand is required")
	ErrVehicleTypeModelRequired = errors.New("vehicle type model is required")
)

var brandCaser = cases.Title(language.Und)

func (s *Server) AddVehicleType(ctx context.Context, brand, model string) (uuid.UUID, error) {
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

	if err := s.VehicleTypeRepo.AddVehicleType(ctx, v); err != nil {
		return uuid.Nil, err
	}

	return v.ID, nil
}

func (s *Server) ListVehicleTypes(ctx context.Context) ([]*entity.VehicleType, error) {
	return s.VehicleTypeRepo.ListVehicleTypes(ctx)
}

func (s *Server) GetVehicleTypeByID(ctx context.Context, id uuid.UUID) (*entity.VehicleType, error) {
	return s.VehicleTypeRepo.GetVehicleTypeByID(ctx, id)
}
