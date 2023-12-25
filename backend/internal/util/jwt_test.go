package util

import (
	"os"
	"testing"
	"time"

	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestJwtEncodeDecode(t *testing.T) {
	pri, err := os.ReadFile("../ed25519_1.key")
	if !assert.Nil(t, err) {
		return
	}
	pub, err := os.ReadFile("../ed25519_1.pub")
	if !assert.Nil(t, err) {
		return
	}

	j := NewJwtBackend(pri, pub, "1")
	claims := &model.UserClaims{
		Id:       1,
		Username: "test",
		Roles:    []string{},
	}
	tokenType := "access"
	token, err := j.Encode(claims, 6*time.Hour, tokenType)
	assert.Nil(t, err)

	assert.NotEmpty(t, token)

	payload, err := j.Decode(token, tokenType)
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, payload, claims)
}
