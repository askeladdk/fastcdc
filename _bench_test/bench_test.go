package bench_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/askeladdk/fastcdc"
	jotfs "github.com/jotfs/fastcdc-go"
	poolporg "github.com/poolpOrg/go-fastcdc"
)

const avgsize = 8 << 10
const datalen = 128 << 20

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (int, error) {
	return fn(p)
}

func randomBytes(num int) []byte {
	rnd := rand.New(rand.NewSource(0))
	rndbytes := make([]byte, num)
	rnd.Read(rndbytes)
	return rndbytes
}

func BenchmarkFastCDC(b *testing.B) {
	rb := randomBytes(datalen)
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	buf := make([]byte, avgsize*16)
	nchunks := 0
	w := writerFunc(func(p []byte) (int, error) {
		nchunks++
		return len(p), nil
	})
	for i := 0; i < b.N; i++ {
		_, _ = fastcdc.CopyBuffer(w, r, buf)
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
}

func BenchmarkJotFS(b *testing.B) {
	rb := randomBytes(datalen)
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	opts := jotfs.Options{
		MinSize:       avgsize / 4,
		AverageSize:   avgsize,
		MaxSize:       avgsize * 8,
		Normalization: 2,
	}
	for i := 0; i < b.N; i++ {
		chunker, _ := jotfs.NewChunker(r, opts)
		for err := error(nil); err == nil; {
			_, err = chunker.Next()
			nchunks++
		}
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
}

func BenchmarkPoolpOrg(b *testing.B) {
	rb := randomBytes(datalen)
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	opts := poolporg.ChunkerOpts{
		MinSize:    avgsize / 4,
		NormalSize: avgsize,
		MaxSize:    avgsize * 8,
	}
	for i := 0; i < b.N; i++ {
		chunker, _ := poolporg.NewChunker(r, &opts)
		for err := error(nil); err == nil; {
			_, err = chunker.Next()
			nchunks++
		}
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
}
