package http

import (
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DownloadWorkOrderHandle returns a PDF work order for a vehicle.
//
// @Summary      Download work order PDF
// @Description  Generates (via create-list-service) and downloads a PDF work order for a vehicle.
// @Tags         work-orders
// @Produce      application/pdf
// @Param        id   path  string  true  "Vehicle ID (uuid)"
// @Success      200  {file}  file
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /vehicles/{id}/work-order [get]
func (h *Handler) DownloadWorkOrderHandle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	vehicleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).SendString(`{"error":"invalid vehicle id"}`)
	}

	pdf, err := h.server.DownloadWorkOrderPDF(c.Context(), vehicleID)
	if err != nil {
		// ошибки пока не маппим (404/500), как договаривались
		return c.Status(500).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	filename := fmt.Sprintf("work-order-%s.pdf", vehicleID.String())

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(pdf)
}
