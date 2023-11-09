package main

import (
	"io"
)

type Buffer struct {
	buf      []byte
	lastRead readOp
	off      int
}
type readOp int8

func (b *Buffer) Position() int {
	return b.off
}
func NewBuffer(buf []byte) *Buffer {

	return &Buffer{buf: buf}
}

const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)

func (b *Buffer) Read(p []byte) (n int, err error) {

	b.lastRead = opInvalid
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return n, nil
}
func (b *Buffer) empty() bool { return len(b.buf) <= b.off }

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (b *Buffer) Next(n int) []byte {
	b.lastRead = opInvalid
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return data
}
func (b *Buffer) Len() int { return len(b.buf) - b.off }
func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
	b.lastRead = opInvalid
}
func (b *Buffer) Bytes() []byte                 { return b.buf[b.off:] }
func (b *Buffer) OffsetBytes(offset int) []byte { return b.buf[offset:] }
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	b.off += len(p)
	return copy(b.buf[b.off:], p), nil
}
