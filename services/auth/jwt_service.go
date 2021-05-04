package auth

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dusansimic/yaas/services"
)

// JWTService is a type that can generate a token from user information or validate a token
// is specified.
type JWTService interface {
	GenerateToken(un string, uid int) (string, error)
	ValidateToken(t string) (*jwt.Token, error)
}

type YaasAuthClaims struct {
	Username string `json:"username"`
	UserID   int    `json:"userid"`
	jwt.StandardClaims
}

type jwtService struct {
	SecretKey string
	Issuer    string
}

// JWTAuthService creates a new jwt authentication service.
func JWTAuthService() JWTService {
	return &jwtService{
		SecretKey: getSecretKey(),
		Issuer:    "yaas",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

// GenerateToken creates a new signed JWT token form specified user account. In case the token
// can be signed, an error is returned.
func (s *jwtService) GenerateToken(un string, uid int) (string, error) {
	claims := &YaasAuthClaims{
		un,
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    s.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.SecretKey))
}

// ValidateToken checks if a specified token is valid. If it's not, or it can't be parsed,
// error is returned.
func (s *jwtService) ValidateToken(t string) (*jwt.Token, error) {
	// Parse specified token and validate it. If token is validated, function will return token
	// structure.
	return jwt.ParseWithClaims(t, &YaasAuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, services.NewServiceError(
				fmt.Errorf("invalid token %s", t.Header["alg"]),
				fmt.Sprintf("invalid token %s", t.Header["alg"]),
			)
		}

		return []byte(s.SecretKey), nil
	})
}
