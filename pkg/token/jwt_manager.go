package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrExpired                 = errors.New("expired token")
	ErrInvalid                 = errors.New("invalid token")
	ErrInvalidClaims           = errors.New("invalid claims")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
)

type JWTClaims struct {
	jwt.StandardClaims
}

type JWTManager struct {
	secret        string
	tokenDuration time.Duration
}

func NewJWTManager(secret string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secret, tokenDuration}
}

func (manager *JWTManager) Generate(id string) (string, error) {
	claims := JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   id,
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secret))
}

func (manager *JWTManager) Verify(accessToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(manager.secret), nil
		},
	)
	if err != nil {
		v, ok := err.(*jwt.ValidationError)
		if ok && v.Errors == jwt.ValidationErrorExpired {
			token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &JWTClaims{})
			if err != nil {
				return nil, ErrInvalid
			}
			return token.Claims.(*JWTClaims), ErrExpired
		}
		return nil, ErrInvalid
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	return claims, nil
}
