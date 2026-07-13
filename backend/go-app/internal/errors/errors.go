package er

import "errors"

var (
	ErrNoCustomer    = errors.New("customer not found")
	ErrNoVehicle     = errors.New("vehicle not found")
	ErrNoVehicleType = errors.New("vehicle type not found")

	ErrGenerateKey = errors.New("sign key generate error")

	ErrCustomerIDRequired = errors.New("customer is required")

	ErrUserAlreadyExists = errors.New("user already exists")
)
