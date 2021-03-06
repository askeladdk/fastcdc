// Package fastcdc implements the fastcdc content-defined chunking (CDC) algorithm.
// CDC is a building block for data deduplication and splits an input stream into
// variable-sized chunks that are likely to be repeated in other, partially similar, inputs.
package fastcdc

//go:generate go run _gen/gear.go -seed 1337 -output gear.go

import (
	"io"
)

const (
	minsize = 2 << 10
	avgsize = 8 << 10
	maxsize = 64 << 10
	bufsize = maxsize << 1
	maskL   = avgsize - 1
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	tail := 0
	head, err := io.ReadFull(src, buf)

	for head > 0 || err == nil {
		i, fp := min(head, tail+minsize), uint16(0)

		for end := min(head, tail+maxsize); i < end; i++ {
			if fp = (fp << 1) + gear[buf[i]]; fp&maskL == 0 {
				break
			}
		}

		if x, err := dst.Write(buf[tail:i]); err != nil {
			return n + int64(x), err
		}

		n, tail = n+int64(i-tail), i

		if unread := head - tail; unread < maxsize {
			if tail > 0 && tail-len(buf) < maxsize {
				copy(buf, buf[tail:])
				tail = 0
			}
			var k int
			if err != io.EOF {
				k, err = io.ReadFull(src, buf[unread:])
			}
			head = tail + unread + k
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
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
	return copyBuffer(dst, src, make([]byte, bufsize))
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer rather than allocating a temporary one.
// If buf is nil, one is allocated; otherwise if it has
// zero length, CopyBuffer panics.
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (n int64, err error) {
	if buf == nil {
		buf = make([]byte, bufsize)
	} else if len(buf) == 0 {
		panic("empty buffer in CopyBuffer")
	}
	return copyBuffer(dst, src, buf)
}
