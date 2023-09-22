package fastcdc

import (
	"hash/fnv"
	"io"
	"math/rand"
	"testing"
	"testing/iotest"
)

func TestCopyIdentical(t *testing.T) {
	const datalen = 1<<20 - 1
	rnd := rand.New(rand.NewSource(0))

	rnd.Seed(0)
	h1 := fnv.New64()
	_, _ = io.Copy(h1, io.LimitReader(rnd, datalen))

	rnd.Seed(0)
	h2 := fnv.New64()
	_, _ = Copy(h2, io.LimitReader(rnd, datalen))

	if h1.Sum64() != h2.Sum64() {
		t.Fatal()
	}
}

func TestCopyErrReader(t *testing.T) {
	_, err := CopyBuffer(io.Discard, iotest.ErrReader(io.ErrClosedPipe), nil)
	if err != io.ErrClosedPipe {
		t.Fatal()
	}
}

func TestCopyRobustness(t *testing.T) {
	count := int64(512<<10 - 1)
	rnd := rand.New(rand.NewSource(0))

	buf := make([]byte, 128<<10)

	for _, testCase := range []struct {
		N string
		R io.Reader
	}{
		{"DataErrReader", iotest.DataErrReader(io.LimitReader(rnd, count))},
		{"OneByteReader", iotest.OneByteReader(io.LimitReader(rnd, count))},
		{"TimeoutReader", iotest.TimeoutReader(io.LimitReader(rnd, count))},
	} {
		t.Run(testCase.N, func(t *testing.T) {
			n, err := CopyBuffer(io.Discard, testCase.R, buf)
			if n != count || err != nil {
				t.Error(n, err)
			}
		})
	}
}

func Benchmark(b *testing.B) {
	for _, x := range []struct {
		Size int
		Name string
	}{
		{1 << 10, "1KB"},
		{4 << 10, "4KB"},
		{16 << 10, "16KB"},
		{64 << 10, "64KB"},
		{256 << 10, "256KB"},
		{1 << 20, "1MB"},
		{4 << 20, "4MB"},
		{16 << 20, "16MB"},
		{64 << 20, "64MB"},
		{256 << 20, "256MB"},
		{1 << 30, "1GB"},
	} {
		x := x
		b.Run(x.Name, func(b *testing.B) {
			buf := make([]byte, bufsize)
			data := make([]byte, x.Size)
			r := bytes.NewReader(data)
			b.ResetTimer()
			b.SetBytes(int64(x.Size))
			for i := 0; i < b.N; i++ {
				r.Reset(data)
				_, _ = CopyBuffer(io.Discard, r, buf)
			}
		})
	}
}
