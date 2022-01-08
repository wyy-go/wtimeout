package wtimeout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBuffer(t *testing.T) {
	pool := &bufferPool{}
	buf := pool.GetBuff()
	assert.NotEqual(t, nil, buf)
	pool.PutBuff(buf)
	buf2 := pool.GetBuff()
	assert.NotEqual(t, nil, buf2)
}
