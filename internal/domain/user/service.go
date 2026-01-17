// internal/domain/user/service.go
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/invoice-app-be/internal/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      Repository
	jwtSecret string
	logger    *logger.Logger // ADD THIS
}

func NewService(repo Repository, jwtSecret string, log *logger.Logger) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    log, // ADD THIS
	}
}

func (s *Service) Register(ctx context.Context, email, password, fullName string) (*User, error) {
	// Check if user exists
	existing, _ := s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("GetByEmail failed", "error", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Error("bcrypt compare failed", "error", err)
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
