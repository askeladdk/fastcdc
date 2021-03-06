package bench_test

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"testing"

	"github.com/askeladdk/fastcdc"
	jotfs "github.com/jotfs/fastcdc-go"
	poolporg "github.com/poolpOrg/go-fastcdc"
	tigerwill90 "github.com/tigerwill90/fastcdc"
)

const (
	minsize = 2 << 10
	avgsize = 8 << 10
	maxsize = 64 << 10
	norm    = 0
	datalen = 128 << 20
)

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (int, error) {
	return fn(p)
}

var rb, _ = io.ReadAll(io.LimitReader(rand.New(rand.NewSource(0)), datalen))

func BenchmarkAskeladdk(b *testing.B) {
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	buf := make([]byte, maxsize<<1)
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

func BenchmarkTigerwill90(b *testing.B) {
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	chunkcounter := func(offset, length uint, chunk []byte) error {
		nchunks++
		return nil
	}
	chunker, _ := tigerwill90.NewChunker(
		context.Background(),
		tigerwill90.WithStreamMode(),
		tigerwill90.WithChunksSize(minsize, avgsize, maxsize),
	)
	for i := 0; i < b.N; i++ {
		_ = chunker.Split(r, chunkcounter)
		_ = chunker.Finalize(chunkcounter)
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
}

func BenchmarkJotFS(b *testing.B) {
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	opts := jotfs.Options{
		MinSize:       minsize,
		AverageSize:   avgsize,
		MaxSize:       maxsize,
		Normalization: norm,
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
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	opts := poolporg.ChunkerOpts{
		MinSize:    minsize,
		NormalSize: avgsize,
		MaxSize:    maxsize,
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
