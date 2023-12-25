package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromFile(t *testing.T) {
	filepath := "../repo/sql/gen.sql"
	queries := LoadFromFile(filepath)
	q1, ok := queries["create"]
	assert.True(t, ok)
	assert.True(t, len(q1) > 0)
	q2, ok := queries["get"]
	assert.True(t, ok)
	assert.True(t, len(q2) > 0)
}
