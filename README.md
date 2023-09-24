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

The package provides `Copy` and `CopyBuffer` functions modeled after the `io` package with identical signatures. The difference is that these Copy functions copy in content-defined chunks instead of fixed-size chunks. Chunks are sized between 32KB and 128KB with an average of about 64KB.

Use `Copy` to copy data from a `io.Reader` to an `io.Writer` in content-defined chunks.

```go
n, err := fastcdc.Copy(w, r)
```

Use `CopyBuffer` to pass a buffer. The buffer size should be 128KB or larger for best results, although it can be smaller. `Copy` allocates a buffer of 256KB. A larger buffer may provide a performance boost by reducing the number of reads.

```go
n, err := fastcdc.CopyBuffer(w, r, make([]byte, 256 << 10))
```

Read the rest of the [documentation on pkg.go.dev](https://godoc.org/github.com/askeladdk/fastcdc). It's easy-peasy!

## Performance

Unscientific benchmarks suggest that this implementation is about 5-10% slower than the fastest implementation (PlakarLabs). As far as I can tell the performance difference is caused by PlakarLabs producing smaller chunks on average which means that it spends less time in the inner loop. Whether that makes it better or worse for deduplication purposes is unclear. However, this implementation makes zero allocations and has the simplest implementation, being less than 100 lines of code including comments.

```sh
% cd _bench_test
% go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: bench_test
cpu: Intel(R) Core(TM) i5-5287U CPU @ 2.90GHz
BenchmarkAskeladdk-4                  20          59166260 ns/op        2268.48 MB/s         54142 avgsz      2479 chunks           52430 B/op          0 allocs/op
BenchmarkTigerwill90-4                15          98254349 ns/op        1366.02 MB/s         66477 avgsz      2019 chunks           17536 B/op          1 allocs/op
BenchmarkJotFS-4                      10         111913617 ns/op        1199.30 MB/s         76828 avgsz      1747 chunks          262256 B/op          2 allocs/op
BenchmarkPlakarLabs-4                 19          53045331 ns/op        2530.25 MB/s         47679 avgsz      2815 chunks          262272 B/op          4 allocs/op
PASS
ok      bench_test      6.947s
```

More benchmarks:

```sh
% go test -run=^$ -bench ^Benchmark$ 
goos: darwin
goarch: amd64
pkg: github.com/askeladdk/fastcdc
cpu: Intel(R) Core(TM) i5-5287U CPU @ 2.90GHz
Benchmark/1KB-4                 19108564                60.02 ns/op     17059.60 MB/s
Benchmark/4KB-4                 12981624                89.85 ns/op     45589.12 MB/s
Benchmark/16KB-4                 3305914               357.1 ns/op      45876.47 MB/s
Benchmark/64KB-4                   41148             29139 ns/op        2249.09 MB/s
Benchmark/256KB-4                  10000            113107 ns/op        2317.66 MB/s
Benchmark/1MB-4                     2394            462801 ns/op        2265.72 MB/s
Benchmark/4MB-4                      636           1805544 ns/op        2323.01 MB/s
Benchmark/16MB-4                     165           7189987 ns/op        2333.41 MB/s
Benchmark/64MB-4                      38          29806177 ns/op        2251.51 MB/s
Benchmark/256MB-4                      9         120255293 ns/op        2232.21 MB/s
Benchmark/1GB-4                        3         479891694 ns/op        2237.47 MB/s
PASS
ok      github.com/askeladdk/fastcdc    28.965s
```

## License

Package fastcdc is released under the terms of the ISC license.
