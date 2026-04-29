package http

import (
	"autofort/internal/entity"
	"autofort/internal/usecase"
	"fmt"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateCustomerRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

type CustomerVehicleDTO struct {
	ID          string `json:"id"`
	VehicleType string `json:"vehicle_type_id"`
	Vin         string `json:"vin"`
	GovNumber   string `json:"gov_number"`
	Mileage     int    `json:"mileage"`
}

type CustomerListItemDTO struct {
	ID          string               `json:"id"`
	FirstName   string               `json:"first_name"`
	LastName    string               `json:"last_name"`
	PhoneNumber string               `json:"phone_number"`
	Vehicles    []CustomerVehicleDTO `json:"vehicles"`
}

type GetCustomerResponse struct {
	Customer CustomerListItemDTO  `json:"customer"`
	Vehicles []CustomerVehicleDTO `json:"vehicles"`
}

type ListCustomersResponse struct {
	Items  []CustomerListItemDTO `json:"items"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Q      string                `json:"q,omitempty"`
}

type UpdateCustomerRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

func fromDTOtoCustomer(req *CreateCustomerRequest) *entity.Customer {
	return &entity.Customer{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}
}

// CreateCustomerHandle creates a new customer
//
// @Summary      Create customer
// @Description  Create a new customer with first name, last name and phone number
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer  body      CreateCustomerRequest  true  "Customer payload"
// @Success      201  {object}  map[string]string  "created customer id"
// @Failure      400  {object}  map[string]string  "validation or bad request error"
// @Router       /customers [post]
func (h *Handler) CreateCustomerHandle(c *fiber.Ctx) error {
	req := &CreateCustomerRequest{}

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	customer := fromDTOtoCustomer(req)

	id, err := h.server.AddCustomer(customer.FirstName, customer.LastName, customer.PhoneNumber)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	return c.Status(201).SendString(fmt.Sprintf(`{"created":"%v"}`, id))

}

// ListCustomersHandle lists customers with vehicles.
//
// @Summary      List customers
// @Description  Search by q (name/lastname/phone). Includes vehicles.
// @Tags         customers
// @Produce      json
// @Param        q       query     string  false  "Search by first name, last name, phone"
// @Param        limit   query     int     false  "Limit (default 20)"
// @Param        offset  query     int     false  "Offset (default 0)"
// @Success      200  {object}  ListCustomersResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /customers [get]
func (h *Handler) ListCustomersHandle(c *fiber.Ctx) error {
	q := c.Query("q", "")

	limit := 20
	if v := c.Query("limit", ""); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return c.Status(400).SendString(`{"error":"invalid limit"}`)
		}
		limit = n
	}

	offset := 0
	if v := c.Query("offset", ""); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return c.Status(400).SendString(`{"error":"invalid offset"}`)
		}
		offset = n
	}

	items, err := h.server.ListCustomers(q, limit, offset)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	resp := ListCustomersResponse{
		Items:  make([]CustomerListItemDTO, 0, len(items)),
		Limit:  limit,
		Offset: offset,
		Q:      q,
	}

	for _, row := range items {
		dto := CustomerListItemDTO{
			ID:          row.Customer.ID.String(),
			FirstName:   row.Customer.FirstName,
			LastName:    row.Customer.LastName,
			PhoneNumber: row.Customer.PhoneNumber,
			Vehicles:    make([]CustomerVehicleDTO, 0, len(row.Vehicles)),
		}

		for _, v := range row.Vehicles {
			dto.Vehicles = append(dto.Vehicles, CustomerVehicleDTO{
				ID:          v.ID.String(),
				VehicleType: v.VehicleType.String(),
				Vin:         v.Vin,
				GovNumber:   v.GovNumber,
				Mileage:     v.Mileage,
			})
		}

		resp.Items = append(resp.Items, dto)
	}

	return c.JSON(resp)
}

// GetCustomerHandle returns customer details with vehicles.
//
// @Summary      Get customer
// @Description  Returns a customer by id with vehicles.
// @Tags         customers
// @Produce      json
// @Param        id   path      string  true  "Customer ID (uuid)"
// @Success      200  {object}  GetCustomerResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /customers/{id} [get]
func (h *Handler) GetCustomerHandle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).SendString(`{"error":"invalid customer id"}`)
	}

	d, err := h.server.GetCustomerDetails(id)
	if err != nil {
		// ошибки пока не маппим (как ты просил), просто отдадим 400/500 потом улучшим
		return c.Status(500).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	customerDTO := CustomerListItemDTO{
		ID:          d.Customer.ID.String(),
		FirstName:   d.Customer.FirstName,
		LastName:    d.Customer.LastName,
		PhoneNumber: d.Customer.PhoneNumber,
		Vehicles:    nil, // не используем тут, чтобы не дублировать
	}

	vehiclesDTO := make([]CustomerVehicleDTO, 0, len(d.Vehicles))
	for _, v := range d.Vehicles {
		vehiclesDTO = append(vehiclesDTO, CustomerVehicleDTO{
			ID:          v.ID.String(),
			VehicleType: v.VehicleType.String(),
			Vin:         v.Vin,
			GovNumber:   v.GovNumber,
			Mileage:     v.Mileage,
		})
	}

	resp := GetCustomerResponse{
		Customer: customerDTO,
		Vehicles: vehiclesDTO,
	}

	return c.JSON(resp)
}

// PatchCustomerHandle updates customer fields.
//
// @Summary      Update customer
// @Description  Partially updates customer fields by id.
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Customer ID (uuid)"
// @Param        payload  body      UpdateCustomerRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /customers/{id} [patch]
func (h *Handler) PatchCustomerHandle(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).SendString(`{"error":"invalid customer id"}`)
	}

	req := new(UpdateCustomerRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	update := usecase.UpdateCustomerInput{
		ID:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}

	if err := h.server.EditCustomer(update); err != nil {
		// ошибки нормально разрулим позже
		return c.Status(500).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	return c.SendStatus(204)
}
