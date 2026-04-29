package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	items, err := h.server.ListVehicleTypes()
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
