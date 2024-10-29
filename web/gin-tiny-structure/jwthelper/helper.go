package jwthelper

import (
	"gdemo/internal/config"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key []byte

func init() {
	if jwtkey, ok := os.LookupEnv("JWTKEY"); !ok {
		panic("Please setup environment variable 'JWTKEY'")
	} else {
		key = []byte(jwtkey)
	}
}

// GenerateToken generates a jwt with a n minutes expiration time.
func GenerateToken() (string, error) {
	config := config.GetConfig()
	now := time.Now()
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":      now.Add(config.JWT.Expiration * time.Minute).Unix(),
			"iat":      now.Unix(),
			"uid":      "42",
			"username": "JD Vance",
		},
	)
	return t.SignedString(key)
}

// ParseToken parses, validates, verifies the tokenString.
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		return key, nil
	})
}
