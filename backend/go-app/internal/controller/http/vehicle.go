package http

import (
	"autofort/internal/usecase"
	"errors"
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddVehicleRequest struct {
	VehicleTypeID string `json:"vehicle_type_id"`
	Vin           string `json:"vin"`
	GovNumber     string `json:"gov_number"`
	Mileage       *int   `json:"mileage,omitempty"`
}

type AddVehicleResponse struct {
	ID string `json:"id"`
}

type UpdateVehicleRequest struct {
	Vin       *string `json:"vin,omitempty"`
	GovNumber *string `json:"gov_number,omitempty"`
	Mileage   *int    `json:"mileage,omitempty"`
}

// PatchVehicleHandle updates vehicle fields.
//
// @Summary      Update vehicle
// @Description  Partially updates vehicle fields by id.
// @Tags         vehicles
// @Accept       json
// @Produce      json
// @Param        id       path      string               true  "Vehicle ID (uuid)"
// @Param        payload  body      UpdateVehicleRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /vehicles/{id} [patch]
func (h *Handler) PatchVehicleHandle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).SendString(`{"error":"invalid vehicle id"}`)
	}

	req := new(UpdateVehicleRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	update := usecase.UpdateVehicleInput{
		ID:        id,
		Vin:       req.Vin,
		GovNumber: req.GovNumber,
		Mileage:   req.Mileage,
	}

	if err := h.server.EditVehicle(update); err != nil {
		return c.Status(500).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	return c.SendStatus(204)
}

// AddVehicleToCustomerHandle adds vehicle to customer.
//
// @Summary      Add vehicle to customer
// @Description  Creates a vehicle for a specific customer (customer id comes from path).
// @Tags         vehicles
// @Accept       json
// @Produce      json
// @Param        id      path      string            true  "Customer ID (uuid)"
// @Param        vehicle body      AddVehicleRequest  true  "Vehicle payload"
// @Success      201     {object}  AddVehicleResponse
// @Failure      400     {object}  map[string]string  "bad request / validation"
// @Failure      404     {object}  map[string]string  "customer or vehicle type not found"
// @Failure      500     {object}  map[string]string  "internal error"
// @Router       /customers/{id}/vehicles [post]
func (h *Handler) AddVehicleToCustomerHandle(c *fiber.Ctx) error {
	// 1) path param
	customerIDStr := c.Params("id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		return c.Status(400).SendString(`{"error":"invalid customer id"}`)
	}

	// 2) body
	req := new(AddVehicleRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	// 3) parse vehicle_type_id
	vtIDStr := strings.TrimSpace(req.VehicleTypeID)
	vtID, err := uuid.Parse(vtIDStr)
	if err != nil {
		return c.Status(400).SendString(`{"error":"invalid vehicle_type_id"}`)
	}

	// 4) call usecase
	vehicleID, err := h.server.AddVehicleToCustomer(
		c.Context(),
		customerID,
		vtID,
		req.Vin,
		req.GovNumber,
		req.Mileage,
	)
	if err != nil {
		// пока минимальный маппинг, позже вынесем в общий error handler
		if errors.Is(err, errors.New("customer not found")) {
			return c.Status(404).SendString(`{"error":"customer not found"}`)
		}
		if errors.Is(err, errors.New("vehicle type not found")) {
			return c.Status(404).SendString(`{"error":"vehicle type not found"}`)
		}
		return c.Status(400).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	return c.Status(201).JSON(AddVehicleResponse{ID: vehicleID.String()})
}
