package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	plainText := "Make America great again"
	hashed, err := HashPassword(plainText)
	assert.Nil(t, err)
	assert.True(t, VerifyPassword(hashed, plainText))
}
