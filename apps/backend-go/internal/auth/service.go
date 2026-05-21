package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("inactive user")
	ErrMissingJWTSecret   = errors.New("jwt secret is required")
)

type Service struct {
	repository Repository
	jwtSecret  []byte
	tokenTTL   time.Duration
	now        func() time.Time
}

func NewService(repository Repository, jwtSecret string, tokenTTL time.Duration) (*Service, error) {
	if strings.TrimSpace(jwtSecret) == "" {
		return nil, ErrMissingJWTSecret
	}
	if tokenTTL <= 0 {
		tokenTTL = time.Duration(defaultTokenTTLSeconds) * time.Second
	}

	return &Service{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
		tokenTTL:   tokenTTL,
		now:        time.Now,
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

	issuedAt := s.now().UTC()
	expiresAt := issuedAt.Add(s.tokenTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"clinic_id": user.ClinicID,
		"role":      user.Role,
		"iat":       issuedAt.Unix(),
		"exp":       expiresAt.Unix(),
	})

	accessToken, err := token.SignedString(s.jwtSecret)
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
