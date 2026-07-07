package usecase

import (
	"autofort/internal/entity"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Postgres interface {
	AddVehicleType(*entity.VehicleType) error

	AddVehicle(*entity.Vehicle) error
	GetVehicle(id uuid.UUID) (*entity.Vehicle, error)
	UpdateVehicle(*entity.Vehicle) error

	AddService(*entity.Service) error
	GetService(id uuid.UUID) (*entity.Service, error)
	SetServiceTime(time float64, serviceId, vehicleId uuid.UUID) error
}

type Server struct {
	CustomerRepo    CustomerRepo
	VehicleTypeRepo VehicleTypeRepo
	VehicleRepo     VehicleRepo
	WorkOrderClient WorkOrderClient
	key             *ecdsa.PrivateKey
}

func NewService(CustomerRepo CustomerRepo, VehicleTypeRepo VehicleTypeRepo, VehicleRepo VehicleRepo, WorkOrderClient WorkOrderClient) *Server {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil
	}
	return &Server{
		CustomerRepo:    CustomerRepo,
		VehicleTypeRepo: VehicleTypeRepo,
		VehicleRepo:     VehicleRepo,
		WorkOrderClient: WorkOrderClient,
		key:             k,
	}
}

func (s *Server) GetKey() *ecdsa.PrivateKey {
	return s.key
}

func (s *Server) GenerateToken(login string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256, // здесь обозначен метод, которым подпишется токен
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   login,
		},
	)
	return token.SignedString(s.key) // здесь токен подписывается
}
