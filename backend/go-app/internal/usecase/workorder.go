package usecase

import (
	"time"

	"github.com/google/uuid"
)

type WorkOrderClient interface {
	CreateWorkOrderPDF(*WorkOrderCreateRequest) (*WorkOrderCreateResponse, error)
}

type WorkOrderCreateRequest struct {
	// Metadata (для шапки документа)
	DocumentID  uuid.UUID `json:"document_id"`
	GeneratedAt time.Time `json:"generated_at"`
	CompanyName string    `json:"company_name,omitempty"` // можно позже
	Title       string    `json:"title,omitempty"`        // "Обходной лист" / "Заказ-наряд"

	Customer CustomerDTO `json:"customer"`
	Vehicle  VehicleDTO  `json:"vehicle"`

	// На будущее (когда появятся orders/services/parts):
	// Services []ServiceLineDTO `json:"services,omitempty"`
	// Parts    []PartLineDTO    `json:"parts,omitempty"`
	// Totals   TotalsDTO        `json:"totals,omitempty"`
}

type CustomerDTO struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
}

type VehicleDTO struct {
	ID        uuid.UUID `json:"id"`
	Vin       string    `json:"vin"`
	GovNumber string    `json:"gov_number"`
	Mileage   int       `json:"mileage"`

	Type VehicleTypeDTO `json:"type"`
}

type VehicleTypeDTO struct {
	ID    uuid.UUID `json:"id"`
	Brand string    `json:"brand"`
	Model string    `json:"model"`
}

type WorkOrderCreateResponse struct {
	PDF []byte
}

func (s *Server) DownloadWorkOrderPDF(vehicleID uuid.UUID) ([]byte, error) {
	// 1) vehicle
	v, err := s.VehicleRepo.GetVehicle(vehicleID)
	if err != nil {
		return nil, err
	}

	// 2) customer
	c, err := s.CustomerRepo.GetCustomer(v.CustomerId)
	if err != nil {
		return nil, err
	}

	// 3) vehicle type
	vt, err := s.VehicleTypeRepo.GetVehicleTypeByID(v.VehicleType)
	if err != nil {
		return nil, err
	}

	req := &WorkOrderCreateRequest{
		DocumentID:  uuid.New(),
		GeneratedAt: time.Now(),
		Title:       "Обходной лист",

		Customer: CustomerDTO{
			ID:          c.ID,
			FirstName:   c.FirstName,
			LastName:    c.LastName,
			PhoneNumber: c.PhoneNumber,
		},
		Vehicle: VehicleDTO{
			ID:        v.ID,
			Vin:       v.Vin,
			GovNumber: v.GovNumber,
			Mileage:   v.Mileage,
			Type: VehicleTypeDTO{
				ID:    vt.ID,
				Brand: vt.Brand,
				Model: vt.Model,
			},
		},
	}

	resp, err := s.WorkOrderClient.CreateWorkOrderPDF(req)
	if err != nil {
		return nil, err
	}

	return resp.PDF, nil
}
