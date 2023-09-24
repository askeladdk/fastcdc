package bench_test

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"testing"

	plakarlabs "github.com/PlakarLabs/go-fastcdc"
	askeladdk "github.com/askeladdk/fastcdc"
	jotfs "github.com/jotfs/fastcdc-go"
	tigerwill90 "github.com/tigerwill90/fastcdc"
)

const (
	minsize = 32 << 10
	avgsize = 64 << 10
	maxsize = 128 << 10
	norm    = 2
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
	buf := make([]byte, 1<<20)
	nchunks := 0
	w := writerFunc(func(p []byte) (int, error) {
		nchunks++
		return len(p), nil
	})
	for i := 0; i < b.N; i++ {
		_, _ = askeladdk.CopyBuffer(w, r, buf)
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
	b.ReportMetric(float64(datalen)/(float64(nchunks)/float64(b.N)), "avgsz")
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
	b.ReportMetric(float64(datalen)/(float64(nchunks)/float64(b.N)), "avgsz")
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
	b.ReportMetric(float64(datalen)/(float64(nchunks)/float64(b.N)), "avgsz")
}

func BenchmarkPlakarLabs(b *testing.B) {
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	opts := plakarlabs.ChunkerOpts{
		MinSize:    minsize,
		NormalSize: avgsize,
		MaxSize:    maxsize,
	}
	for i := 0; i < b.N; i++ {
		chunker, _ := plakarlabs.NewChunker(r, &opts)
		for err := error(nil); err == nil; {
			_, err = chunker.Next()
			nchunks++
		}
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
	b.ReportMetric(float64(datalen)/(float64(nchunks)/float64(b.N)), "avgsz")
}
