package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnowflake_New(t *testing.T) {
	id := Generate()
	t.Log(id)
	assert.Equal(t, len(id), 19)
}
