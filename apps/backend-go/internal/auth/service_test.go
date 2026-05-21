package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepository struct {
	user User
	err  error
}

func (r fakeRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	if r.err != nil {
		return User{}, r.err
	}
	return r.user, nil
}

func TestLoginReturnsTokenForValidCredentials(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("valid-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	service, err := NewService(fakeRepository{user: User{
		ID:           "user-id",
		ClinicID:     "clinic-id",
		FullName:     "Ana Gomez",
		Email:        "assistant@clinic.example",
		PasswordHash: string(hash),
		Role:         RoleAssistant,
		IsActive:     true,
	}}, "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	service.now = func() time.Time { return time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC) }

	res, err := service.Login(context.Background(), LoginRequest{Email: "assistant@clinic.example", Password: "valid-password"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if res.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if res.TokenType != TokenTypeBearer {
		t.Fatalf("expected token type %s, got %s", TokenTypeBearer, res.TokenType)
	}
	if res.ExpiresIn != 3600 {
		t.Fatalf("expected expires_in 3600, got %d", res.ExpiresIn)
	}
	if res.User.ID != "user-id" || res.User.ClinicID != "clinic-id" || res.User.Email != "assistant@clinic.example" {
		t.Fatalf("unexpected user response: %#v", res.User)
	}

	parsed, err := jwt.Parse(res.AccessToken, func(token *jwt.Token) (any, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	claims := parsed.Claims.(jwt.MapClaims)
	if claims["sub"] != "user-id" || claims["clinic_id"] != "clinic-id" || claims["role"] != RoleAssistant {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("valid-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	service, err := NewService(fakeRepository{user: User{PasswordHash: string(hash), IsActive: true}}, "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, err = service.Login(context.Background(), LoginRequest{Email: "assistant@clinic.example", Password: "wrong"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginRejectsInactiveUser(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("valid-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	service, err := NewService(fakeRepository{user: User{PasswordHash: string(hash), IsActive: false}}, "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, err = service.Login(context.Background(), LoginRequest{Email: "assistant@clinic.example", Password: "valid-password"})
	if !errors.Is(err, ErrInactiveUser) {
		t.Fatalf("expected ErrInactiveUser, got %v", err)
	}
}

func TestNewServiceRequiresJWTSecret(t *testing.T) {
	_, err := NewService(fakeRepository{}, "", time.Hour)
	if !errors.Is(err, ErrMissingJWTSecret) {
		t.Fatalf("expected ErrMissingJWTSecret, got %v", err)
	}
}
