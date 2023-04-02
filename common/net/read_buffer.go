package net

import (
	"bufio"
	"net"
)
/* 包装标准库net.Conn, 组合bufio.Reader的功能*/
type BufReader struct {
	r *bufio.Reader
	net.Conn
}

func NewReader(c net.Conn) (b *BufReader) {
	if b, ok := c.(*BufReader); ok {
		return b
	}
	
	b = &BufReader{
		r:    bufio.NewReader(c),
		Conn: c,
	}
	return
}

func (b *BufReader) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b *BufReader) Read(p []byte) (n int, err error) {
	return b.r.Read(p)
}

// Reader returns the internal bufio.Reader.
func (b *BufReader) Reader() *bufio.Reader {
	return b.r
}
