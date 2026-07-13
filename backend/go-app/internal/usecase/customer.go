package usecase

import (
	"autofort/internal/entity"

	"github.com/google/uuid"
)

type CustomerRepo interface {
	AddCustomer(*entity.Customer) error
	GetCustomer(id uuid.UUID) (*entity.Customer, error)
	UpdateCustomer(*entity.Customer) error
	SearchCustomers(q string, limit, offset int) ([]*entity.Customer, error)
}

type CustomerWithVehicles struct {
	Customer *entity.Customer
	Vehicles []*entity.Vehicle
}

type UpdateCustomerInput struct {
	ID          uuid.UUID
	FirstName   *string
	LastName    *string
	PhoneNumber *string
}
type CustomerDetails struct {
	Customer *entity.Customer
	Vehicles []*entity.Vehicle
}

func (s *Server) GetCustomerDetails(id uuid.UUID) (*CustomerDetails, error) {
	c, err := s.CustomerRepo.GetCustomer(id)
	if err != nil {
		return nil, err
	}

	vehicles, err := s.VehicleRepo.ListVehiclesByCustomerIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	return &CustomerDetails{
		Customer: c,
		Vehicles: vehicles,
	}, nil
}

func (s *Server) AddCustomer(f, l, phone string) (uuid.UUID, error) {
	c := &entity.Customer{
		ID:          uuid.New(),
		FirstName:   f,
		LastName:    l,
		PhoneNumber: phone,
	}

	if err := c.Validate(); err != nil {
		return uuid.Nil, err
	}

	if err := s.CustomerRepo.AddCustomer(c); err != nil {
		return uuid.Nil, err
	}

	return c.ID, nil
}

func (s *Server) EditCustomer(update UpdateCustomerInput) error {
	c, err := s.CustomerRepo.GetCustomer(update.ID)
	if err != nil {
		return err
	}

	if update.FirstName != nil {
		c.FirstName = *update.FirstName
	}

	if update.LastName != nil {
		c.LastName = *update.LastName
	}

	if update.PhoneNumber != nil {
		c.PhoneNumber = *update.PhoneNumber
	}

	if err = c.Validate(); err != nil {
		return err
	}

	return s.CustomerRepo.UpdateCustomer(c)
}

func (s *Server) GetCustomerByID(ID string) (*entity.Customer, error) {
	id, err := uuid.Parse(ID)
	if err != nil {
		return nil, err
	}

	return s.CustomerRepo.GetCustomer(id)
}

func (s *Server) ListCustomers(q string, limit, offset int) ([]CustomerWithVehicles, error) {
	customers, err := s.CustomerRepo.SearchCustomers(q, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(customers) == 0 {
		return []CustomerWithVehicles{}, nil
	}

	ids := make([]uuid.UUID, 0, len(customers))
	for _, c := range customers {
		ids = append(ids, c.ID)
	}

	vehicles, err := s.VehicleRepo.ListVehiclesByCustomerIDs(ids)
	if err != nil {
		return nil, err
	}

	byCustomer := make(map[uuid.UUID][]*entity.Vehicle, len(customers))
	for _, v := range vehicles {
		byCustomer[v.CustomerId] = append(byCustomer[v.CustomerId], v)
	}

	out := make([]CustomerWithVehicles, 0, len(customers))
	for _, c := range customers {
		out = append(out, CustomerWithVehicles{
			Customer: c,
			Vehicles: byCustomer[c.ID],
		})
	}

	return out, nil
}
