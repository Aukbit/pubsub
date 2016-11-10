package pubsub

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestHello(t *testing.T) {
	m := hello()
	assert.Equal(t, "hello", m)
}
