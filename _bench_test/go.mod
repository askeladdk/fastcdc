module bench_test

go 1.21

toolchain go1.21.1

require (
	github.com/PlakarLabs/go-fastcdc v0.5.0
	github.com/askeladdk/fastcdc v0.0.0
	github.com/jotfs/fastcdc-go v0.2.0
	github.com/tigerwill90/fastcdc v1.2.2
)

replace github.com/askeladdk/fastcdc => ../
