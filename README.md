# fastcdc - fast content-defined chunking in Go

[![GoDoc](https://godoc.org/github.com/askeladdk/fastcdc?status.png)](https://godoc.org/github.com/askeladdk/fastcdc)
[![Go Report Card](https://goreportcard.com/badge/github.com/askeladdk/fastcdc)](https://goreportcard.com/report/github.com/askeladdk/fastcdc)

## Overview

Package fastcdc implements the [fastcdc](https://www.usenix.org/system/files/conference/atc16/atc16-paper-xia.pdf) content-defined chunking (CDC) algorithm. CDC is a building block for data deduplication and splits an input stream into variable-sized chunks that are likely to be repeated in other, partially similar, inputs.

## Install

```
go get -u github.com/askeladdk/fastcdc
```

## Quickstart

The package provides `Copy` and `CopyBuffer` functions modeled after the `io` package with identical signatures. The difference is that these Copy functions copy in content-defined chunks instead of fixed-size chunks. Chunks are sized between 8KB and 32KB with an average of about 16KB.

Use `Copy` to copy data from a `io.Reader` to an `io.Writer` in content-defined chunks.

```go
n, err := fastcdc.Copy(w, r)
```

Use `CopyBuffer` to pass a buffer. The buffer size should be 64KB or larger for best results, although it can be smaller. `Copy` allocates a buffer of 64KB. A larger buffer may provide a performance boost by reducing the number of reads.

```go
n, err := fastcdc.CopyBuffer(w, r, make([]byte, 256 << 10))
```

Use `Chunker` to customize the parameters:

```go
chunker := fastcdc.Chunker {
    MinSize: 1 << 20,
    AvgSize: 2 << 20,
    MaxSize: 4 << 20,
    Norm:    2,
}

buf := make([]byte, 2*chunker.MaxSize)
n, err := chunker.CopyBuffer(dst, src, buf)
```

Read the rest of the [documentation on pkg.go.dev](https://godoc.org/github.com/askeladdk/fastcdc). It's easy-peasy!

## Performance

Unscientific benchmarks suggest that this implementation is about as fast as Tigerwill90 but produces larger chunks. This is due to Tigerwill90's slightly different fingerprint calculation (they shift right instead of left). PlakarLabs has much higher performance but this is because it produces smaller chunks, meaning that it spends less time in the inner loop.

Unlike the others, this implementation makes zero allocations and only has the fewest lines of code.

```sh
% cd _bench_test
% go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: bench_test
cpu: Intel(R) Core(TM) i5-5287U CPU @ 2.90GHz
BenchmarkAskeladdk-4                  14          78664269 ns/op        1706.21 MB/s       2485513 avgsz        54.00 chunks       599188 B/op          0 allocs/op
BenchmarkTigerwill90-4                13          77380696 ns/op        1734.51 MB/s       2064888 avgsz        65.00 chunks       645339 B/op          1 allocs/op
BenchmarkJotFS-4                      10         103483790 ns/op        1296.99 MB/s       2396745 avgsz        56.00 chunks      8388720 B/op          2 allocs/op
BenchmarkPlakarLabs-4                 31          36523149 ns/op        3674.87 MB/s       1065220 avgsz       126.0 chunks       8388736 B/op          4 allocs/op
PASS
ok      bench_test      5.136s
```

More unscientific benchmarks:

```sh
% go test -run=^$ -bench ^Benchmark$ 
goos: darwin
goarch: amd64
pkg: github.com/askeladdk/fastcdc
cpu: Intel(R) Core(TM) i5-5287U CPU @ 2.90GHz
Benchmark/1KB-4                  8513276               120.5 ns/op      8497.58 MB/s
Benchmark/4KB-4                  6978042               153.9 ns/op      26619.10 MB/s
Benchmark/16KB-4                  166795              7117 ns/op        2302.14 MB/s
Benchmark/64KB-4                   53578             22183 ns/op        2954.29 MB/s
Benchmark/256KB-4                   9573            122433 ns/op        2141.11 MB/s
Benchmark/1MB-4                     2134            521845 ns/op        2009.36 MB/s
Benchmark/4MB-4                      534           2116966 ns/op        1981.28 MB/s
Benchmark/16MB-4                     140           8525421 ns/op        1967.90 MB/s
Benchmark/64MB-4                      33          34171293 ns/op        1963.90 MB/s
Benchmark/256MB-4                      8         135296222 ns/op        1984.06 MB/s
Benchmark/1GB-4                        2         548831781 ns/op        1956.41 MB/s
PASS
ok      github.com/askeladdk/fastcdc    22.673s
```

## License

Package fastcdc is released under the terms of the ISC license.
