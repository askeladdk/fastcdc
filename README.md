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

The package provides `Copy` and `CopyBuffer` functions modeled after the `io` package with identical signatures. The difference is that these Copy functions copy in content-defined chunks instead of fixed-size chunks.

Use `Copy` to copy data from a `io.Reader` to an `io.Writer` in content-defined chunks.

```go
n, err := fastcdc.Copy(w, r)
```

Use `CopyBuffer` to pass a buffer. The buffer size should be 64KB or larger for best results, although it can be smaller. `Copy` allocates a buffer of 128KB.

```go
n, err := fastcdc.CopyBuffer(w, r, make([]byte, 128 << 10))
```

Unlike other implementations it is not possible to tweak the parameters. This is not needed because there is a sweet spot of practical chunk sizes that enables efficient deduplication: Too small reduces performance due to overhead and too high reduces deduplication due to overly coarse chunks. Hence, chunks are sized between 2KB and 64KB with an average of 8KB. The final chunk can be smaller than 2KB.

Read the rest of the [documentation on pkg.go.dev](https://godoc.org/github.com/askeladdk/fastcdc). It's easy-peasy!

## Performance

Unscientific benchmarks suggest that our implementation is 20% faster than the next best implementation and is the only one that makes zero allocations. Our implementation is also much simpler than the others, being less than 100 lines of code including comments. The number of generated chunks is roughly the same.

```sh
% cd _bench_test
% go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: bench_test
cpu: Intel(R) Core(TM) i5-5287U CPU @ 2.90GHz
BenchmarkAskeladdk-4     	      14	  82502821 ns/op	1626.83 MB/s	     14348 chunks	    9364 B/op	       0 allocs/op
BenchmarkTigerwill90-4   	      12	  99786487 ns/op	1345.05 MB/s	     16027 chunks	   10985 B/op	       1 allocs/op
BenchmarkJotFS-4         	      10	 107251454 ns/op	1251.43 MB/s	     14651 chunks	  131184 B/op	       2 allocs/op
BenchmarkPoolpOrg-4      	       4	 250803846 ns/op	 535.15 MB/s	     14328 chunks	144083848 B/op	   42990 allocs/op
PASS
ok  	bench_test	8.117s
```

## License

Package fastcdc is released under the terms of the ISC license.
