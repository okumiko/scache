package bufferpool

import (
	"bytes"
	"sync"
)

type bufferPool struct {
	Pool sync.Pool
}

func NewBufferPool() *bufferPool {
	return &bufferPool{
		Pool: sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}
}
