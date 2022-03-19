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
