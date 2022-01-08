package wtimeout

import (
	"bytes"
	"sync"
)

const buffSize = 10 * 1024

var defaultBufPool bufferPool

// BufferPool is Pool of *bytes.Buffer
type bufferPool struct {
	pool sync.Pool
}

// GetBuff a bytes.Buffer pointer
func (p *bufferPool) GetBuff() *bytes.Buffer {
	buf := p.pool.Get()
	if buf == nil {
		bs := make([]byte, 0, buffSize)
		return bytes.NewBuffer(bs)
	}
	return buf.(*bytes.Buffer)
}

// PutBuff a bytes.Buffer pointer to BufferPool
func (p *bufferPool) PutBuff(buf *bytes.Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}
