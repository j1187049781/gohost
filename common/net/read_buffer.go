package net

import (
	"bufio"
	"net"
)

type ReadBufCONN struct {
	r *bufio.Reader
	net.Conn
}

func NewReadBufCONN(c net.Conn)(b *ReadBufCONN){
	b = &ReadBufCONN{
		r: bufio.NewReader(c),
		Conn: c,
	}
	return
}
func (b *ReadBufCONN)Peek(n int) ([]byte, error){
	return b.r.Peek(n)
}
func (b *ReadBufCONN)Read(p []byte) (n int, err error){
	return b.r.Read(p)
}
// Reader returns the internal bufio.Reader.
func (b *ReadBufCONN) Reader() *bufio.Reader {
	return b.r
}
