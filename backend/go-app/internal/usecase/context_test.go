package usecase

import (
	"context"
	"errors"
	"testing"

	"autofort/internal/entity"

	"github.com/google/uuid"
)

type customerContextRepo struct {
	CustomerRepo
	get func(context.Context, uuid.UUID) (*entity.Customer, error)
}

func (r customerContextRepo) GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	return r.get(ctx, id)
}

type vehicleContextRepo struct {
	VehicleRepo
	get  func(context.Context, uuid.UUID) (*entity.Vehicle, error)
	list func(context.Context, []uuid.UUID) ([]*entity.Vehicle, error)
}

func (r vehicleContextRepo) GetVehicle(ctx context.Context, id uuid.UUID) (*entity.Vehicle, error) {
	return r.get(ctx, id)
}

func (r vehicleContextRepo) ListVehiclesByCustomerIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Vehicle, error) {
	return r.list(ctx, ids)
}

func TestGetCustomerDetailsStopsAfterRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	id := uuid.New()
	listed := false
	s := &Server{
		CustomerRepo: customerContextRepo{get: func(received context.Context, receivedID uuid.UUID) (*entity.Customer, error) {
			if received != ctx || receivedID != id {
				t.Fatal("customer lookup did not receive the request context and ID")
			}
			cancel()
			return &entity.Customer{ID: id}, nil
		}},
		VehicleRepo: vehicleContextRepo{list: func(received context.Context, ids []uuid.UUID) ([]*entity.Vehicle, error) {
			listed = true
			if len(ids) != 1 || ids[0] != id {
				t.Fatalf("unexpected customer IDs: %v", ids)
			}
			return nil, received.Err()
		}},
	}

	result, err := s.GetCustomerDetails(ctx, id)
	if !listed || result != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation from vehicle lookup, got result=%v, err=%v, listed=%v", result, err, listed)
	}
}

func TestDownloadWorkOrderPDFStopsAfterRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vehicleID, customerID := uuid.New(), uuid.New()
	lookedUpCustomer := false
	s := &Server{
		VehicleRepo: vehicleContextRepo{get: func(received context.Context, id uuid.UUID) (*entity.Vehicle, error) {
			if received != ctx || id != vehicleID {
				t.Fatal("vehicle lookup did not receive the request context and ID")
			}
			cancel()
			return &entity.Vehicle{ID: vehicleID, CustomerId: customerID}, nil
		}},
		CustomerRepo: customerContextRepo{get: func(received context.Context, id uuid.UUID) (*entity.Customer, error) {
			lookedUpCustomer = true
			if id != customerID {
				t.Fatalf("unexpected customer ID: %v", id)
			}
			return nil, received.Err()
		}},
	}

	pdf, err := s.DownloadWorkOrderPDF(ctx, vehicleID)
	if !lookedUpCustomer || pdf != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation before rendering, got pdf=%v, err=%v, lookedUpCustomer=%v", pdf, err, lookedUpCustomer)
	}
}
