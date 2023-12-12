package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// install go jwt library
// go get -u github.com/golang-jwt/jwt/v5
// JWTManager is a JSON web token manager
type JWTManager struct {
	sercetKey     string
	tokenDuration time.Duration
}

// UserClaims is a custom JWT claims that contains some user's information
type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

// NewJWTManager returns a new JWT manager
func NewJWTManager(sercetKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{sercetKey, tokenDuration}
}

func (manager *JWTManager) Generate(user *User) (string, error) {

	nowTime := time.Now()
	expireTime := nowTime.Add(manager.tokenDuration)

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
		Username: user.Username,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.sercetKey))
}

// Verify verifies the access token string and return a user claims if the token is valid
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.sercetKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims := token.Claims.(*UserClaims)

	return claims, nil
}
