package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 24 * time.Hour
)

type tokens struct {
	// for one session only
	access string
	// long time creedntials
	refreash string
}

type UserRepo interface {
	GetPassHash(ctx context.Context, userID *uuid.UUID, nickname string, hash *string) error
}

func (s *Server) SignIn(ctx context.Context, nickname, password string) (*tokens, error) {
	var (
		userID uuid.UUID
		hash   string
	)

	err := s.UserRepo.GetPassHash(ctx, &userID, nickname, &hash)
	if err != nil {
		return nil, fmt.Errorf("signIn error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("wrong password: %w", err)
	}

	access, err := s.generateToken(userID.String(), "access", accessTTL)
	if err != nil {
		return nil, fmt.Errorf("generate access error: %w", err)
	}

	refresh, err := s.generateToken(userID.String(), "refresh", refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("generate refresh error: %w", err)
	}

	return &tokens{
		access:   access,
		refreash: refresh,
	}, nil
}

func (s *Server) generateToken(userID, jti string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   userID,
			ID:        jti,
		},
	)
	return token.SignedString(s.key)
}
