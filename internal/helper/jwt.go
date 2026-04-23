package helper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken  = errors.New("invalid or expired token")
	ErrInvalidClaims = errors.New("invalid token claims")
)

func AccessSecretKey() []byte {
	key := os.Getenv("JWT_ACCESS_SECRET")
	if key == "" {
		panic("JWT_ACCESS_SECRET is not set")
	}
	return []byte(key)
}

func RefreshSecretKey() []byte {
	key := os.Getenv("JWT_REFRESH_SECRET")
	if key == "" {
		panic("JWT_REFRESH_SECRET is not set")
	}
	return []byte(key)
}

type JwtPayload struct {
	ID    uuid.UUID
	Email string
}

func GenerateAccessToken(payload JwtPayload) (string, error) {
	claims := jwt.MapClaims{
		"sub":   payload.ID,
		"email": payload.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessSecretKey())
}

func GenerateRefreshToken(payload JwtPayload) (string, error) {
	claims := jwt.MapClaims{
		"sub":   payload.ID,
		"email": payload.Email,
		"exp":   time.Now().Add(time.Hour * 7 * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(RefreshSecretKey())
}

func VerifyAccessToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return AccessSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return &claims, nil
}

func VerifyRefreshToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return RefreshSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return &claims, nil
}
