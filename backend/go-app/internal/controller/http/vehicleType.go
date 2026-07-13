package http

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VehicleTypeResponse struct {
	ID    string `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
}

type ListVehicleTypesResponse struct {
	Items []VehicleTypeResponse `json:"items"`
}

// GetVehicleTypesHandle lists available vehicle types.
//
// @Summary      List vehicle types
// @Description  Returns vehicle types (brand + model) for dropdowns.
// @Tags         vehicle-types
// @Accept       json
// @Produce      json
// @Success      200  {object}  ListVehicleTypesResponse
// @Failure      500  {object}  map[string]string  "internal error"
// @Router       /vehicle-types [get]
func (h *Handler) GetVehicleTypesHandle(c *fiber.Ctx) error {
	items, err := h.server.ListVehicleTypes(c.Context())
	if err != nil {
		return c.Status(404).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	resp := ListVehicleTypesResponse{
		Items: make([]VehicleTypeResponse, 0, len(items)),
	}

	for _, vt := range items {
		resp.Items = append(resp.Items, VehicleTypeResponse{
			ID:    vt.ID.String(),
			Brand: vt.Brand,
			Model: vt.Model,
		})
	}

	return c.Status(200).JSON(resp)
}

func (h *Handler) GetVehicleTypeByIDHandle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).SendString(`{"error": "invalid vehicle id"}`)
	}

	vp, err := h.server.GetVehicleTypeByID(c.Context(), id)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	r := &VehicleTypeResponse{
		ID:    vp.ID.String(),
		Brand: vp.Brand,
		Model: vp.Model,
	}

	resp, err := json.Marshal(r)
	if err != nil {
		return c.Status(400).SendString(`{"error": "invalid marshalling"}`)
	}

	return c.Status(201).JSON(resp)
}
