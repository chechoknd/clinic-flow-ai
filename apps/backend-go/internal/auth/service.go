package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("inactive user")
	ErrMissingJWTSecret   = errors.New("jwt secret is required")
)

type Service struct {
	repository   Repository
	tokenManager *TokenManager
	tokenTTL     time.Duration
	now          func() time.Time
}

func NewService(repository Repository, jwtSecret string, tokenTTL time.Duration) (*Service, error) {
	tokenManager, err := NewTokenManager(jwtSecret)
	if err != nil {
		return nil, err
	}
	if tokenTTL <= 0 {
		tokenTTL = time.Duration(defaultTokenTTLSeconds) * time.Second
	}

	return &Service{
		repository:   repository,
		tokenManager: tokenManager,
		tokenTTL:     tokenTTL,
		now:          time.Now,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	email := strings.TrimSpace(req.Email)
	password := req.Password
	if email == "" || password == "" {
		return LoginResponse{}, ErrInvalidCredentials
	}

	user, err := s.repository.FindByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return LoginResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResponse{}, err
	}
	if !user.IsActive {
		return LoginResponse{}, ErrInactiveUser
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.tokenManager.Sign(user, s.now(), s.tokenTTL)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: accessToken,
		TokenType:   TokenTypeBearer,
		ExpiresIn:   int64(s.tokenTTL.Seconds()),
		User: UserResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Role:     user.Role,
			ClinicID: user.ClinicID,
		},
	}, nil
}
