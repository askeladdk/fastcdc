//go:build ignore

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	seed := flag.Int64("seed", 0, "random seed")
	outp := flag.String("output", "gear.go", "output filename")
	flag.Parse()

	_ = os.Remove(*outp)

	f, _ := os.OpenFile(*outp, os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	rnd := rand.New(rand.NewSource(*seed))
	fmt.Fprintln(f, "// Code generated with _gen/gear.go. DO NOT EDIT.\n")
	fmt.Fprintln(f, "package fastcdc\n")
	fmt.Fprintf(f, "// gear is generated from rand.New(rand.NewSource(%d)).\n", *seed)
	fmt.Fprintln(f, "var gear [256]uint64 = [...]uint64{")
	for i := 0; i < 64; i++ {
		fmt.Fprintf(f, "\t")
		for j := 0; j < 3; j++ {
			fmt.Fprintf(f, "0x%016x, ", rnd.Uint64())
		}
		fmt.Fprintf(f, "0x%016x,\n", rnd.Uint64())
	}
	fmt.Fprintln(f, "}")
}
