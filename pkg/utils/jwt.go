package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	tokenSigningKey   []byte
	refreshSigningKey []byte
	tokenDuration     time.Duration
	refreshDuration   time.Duration
}

type Claims struct {
	UserId string
	jwt.RegisteredClaims
}

func NewJwt(
	tokenSecret string,
	refreshSecret string,
	tokenDurationInSec int,
	refreshDurationInSec int,
) *Jwt {
	return &Jwt{
		tokenSigningKey:   []byte(tokenSecret),
		refreshSigningKey: []byte(refreshSecret),
		tokenDuration:     time.Duration(tokenDurationInSec) * time.Second,
		refreshDuration:   time.Duration(refreshDurationInSec) * time.Second,
	}
}

func (j *Jwt) GenerateAccessToken(userId string) (string, error) {
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return generateTokenWithKey(j.tokenSigningKey, claims)
}

func (j *Jwt) GenerateRefreshToken(userId string) (string, error) {
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return generateTokenWithKey(j.refreshSigningKey, claims)
}

func (j *Jwt) ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	return parseTokenWithKey(j.tokenSigningKey, tokenString)
}

func (j *Jwt) ParseRefreshToken(tokenString string) (*jwt.Token, *Claims, error) {
	return parseTokenWithKey(j.refreshSigningKey, tokenString)
}

func (j *Jwt) ValidateToken(tokenString string) (bool, *Claims) {
	token, claims, err := j.ParseToken(tokenString)
	if err != nil {
		return false, claims
	}

	return token.Valid, claims
}

func (j *Jwt) ValidateRefreshToken(tokenString string) (bool, *Claims) {
	token, claims, err := j.ParseRefreshToken(tokenString)
	if err != nil {
		return false, nil
	}

	return token.Valid, claims
}

func generateTokenWithKey(signedKey []byte, claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signedKey)
}

func parseTokenWithKey(signedKey []byte, tokenString string) (*jwt.Token, *Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return signedKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, jwt.ErrSignatureInvalid
}
