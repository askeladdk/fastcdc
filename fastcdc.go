// Package fastcdc implements the fastcdc content-defined chunking (CDC) algorithm.
// CDC is a building block for data deduplication and splits an input stream into
// variable-sized chunks that are likely to be repeated in other, partially similar, inputs.
package fastcdc

//go:generate go run _gen/gear.go -seed 1337 -output gear.go

import (
	"io"
	"math"
)

// Chunker is a configurable content defined chunker.
type Chunker struct {
	// MinSize is the minimum chunk size in bytes.
	MinSize int
	// AvgSize is the average chunk size in bytes.
	AvgSize int
	// MaxSize is the maximum chunk size in bytes.
	MaxSize int
	// Norm is the normalization factor. Set to zero to disable normalization.
	Norm int
}

// Copy copies from src to dst in content-defined chunk sizes.
// A successful Copy returns err == nil.
func (c Chunker) Copy(dst io.Writer, src io.Reader) (n int64, err error) {
	return c.copyBuffer(dst, src, make([]byte, c.MaxSize<<1))
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer rather than allocating a temporary one.
func (c Chunker) CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	if buf == nil {
		buf = make([]byte, c.MaxSize<<1)
	} else if len(buf) == 0 {
		panic("fastcdc: empty buffer in CopyBuffer")
	}
	return c.copyBuffer(dst, src, buf)
}

func (c Chunker) copyBuffer(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	bits := int(math.Floor(math.Log2(float64(c.AvgSize))))
	maskS := uint64(1)<<max(0, min(bits+c.Norm, 64)) - 1
	maskL := uint64(1)<<max(0, min(bits-c.Norm, 64)) - 1
	gear := gear // speeds up the inner loop

	tail := 0
	head, err := io.ReadFull(src, buf)

	for head > 0 || err == nil {
		i := min(head, tail+c.MinSize)
		fp := uint64(0)

		for m, j := maskS, min(head, tail+c.AvgSize); i < j; i++ {
			if fp = fp<<1 + gear[buf[i]]; fp&m == 0 {
				goto emitchunk
			}
		}

		for m, j := maskL, min(head, tail+c.MaxSize); i < j; i++ {
			if fp = fp<<1 + gear[buf[i]]; fp&m == 0 {
				break
			}
		}

	emitchunk:
		if x, err := dst.Write(buf[tail:i]); err != nil {
			return n + int64(x), err
		}

		n, tail = n+int64(i-tail), i

		if unread := head - tail; unread < c.MaxSize {
			copy(buf, buf[tail:head])
			var k int
			if err != io.EOF {
				k, err = io.ReadFull(src, buf[unread:])
			}
			tail, head = 0, unread+k
		}
	}

	if err == io.EOF {
		err = nil
	}

	return n, err
}

var defaultChunker = Chunker{
	MinSize: 8 << 10,
	AvgSize: 16 << 10,
	MaxSize: 32 << 10,
	Norm:    2,
}

// Copy copies from src to dst in content-defined chunk sizes,
// as opposed to io.Copy which copies in fixed-sized chunks.
//
// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
func Copy(dst io.Writer, src io.Reader) (n int64, err error) {
	return defaultChunker.Copy(dst, src)
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer rather than allocating a temporary one.
// If buf is nil, one is allocated; otherwise if it has
// zero length, CopyBuffer panics.
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	return defaultChunker.CopyBuffer(dst, src, buf)
}

// DefaultChunker returns the chunker used by [Copy] and [CopyBuffer].
func DefaultChunker() Chunker {
	return defaultChunker
}
