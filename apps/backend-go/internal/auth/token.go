package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Claims struct {
	UserID   string
	ClinicID string
	Role     string
}

type TokenManager struct {
	secret []byte
}

func NewTokenManager(jwtSecret string) (*TokenManager, error) {
	if strings.TrimSpace(jwtSecret) == "" {
		return nil, ErrMissingJWTSecret
	}
	return &TokenManager{secret: []byte(jwtSecret)}, nil
}

func (m *TokenManager) Sign(user User, issuedAt time.Time, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = time.Duration(defaultTokenTTLSeconds) * time.Second
	}
	issuedAt = issuedAt.UTC()
	expiresAt := issuedAt.Add(ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"clinic_id": user.ClinicID,
		"role":      user.Role,
		"iat":       issuedAt.Unix(),
		"exp":       expiresAt.Unix(),
	})
	return token.SignedString(m.secret)
}

func (m *TokenManager) Parse(accessToken string) (Claims, error) {
	parsed, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return Claims{}, ErrExpiredToken
	}
	if err != nil || !parsed.Valid {
		return Claims{}, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, ErrInvalidToken
	}
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return Claims{}, ErrInvalidToken
	}
	clinicID, ok := claims["clinic_id"].(string)
	if !ok || clinicID == "" {
		return Claims{}, ErrInvalidToken
	}
	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return Claims{}, ErrInvalidToken
	}

	return Claims{UserID: userID, ClinicID: clinicID, Role: role}, nil
}
