package jwthelper

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	tokenString, err := GenerateToken()
	assert.Nil(t, err)
	token, err := ParseToken(tokenString)
	assert.Nil(t, err)
	assert.Equal(t, tokenString, token.Raw)
	assert.True(t, token.Valid)
	claims := *(token.Claims.(*jwt.MapClaims))
	assert.Equal(t, "42", claims["uid"])
}
