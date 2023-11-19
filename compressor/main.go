package main

import (
	"errors"
	"flag"
	"fmt"
	"nkwatra/compressor/cmp"
	"os"
)

func main() {
	defer func() {
		code := recover()
		if val, isError := code.(error); isError {
			fmt.Println(val)
			os.Exit(1)
		}
	}()
	outPtr := flag.String("out", "", "Path to output file")
	cPtr := flag.Bool("e", false, "To encode/compress a file")
	dPtr := flag.Bool("d", false, "To decode/expand a file")
	flag.Parse()
	if *outPtr == "" {
		panic(errors.New("compressor: output path is required"))
	}
	if !(*cPtr) && !(*dPtr) {
		panic(errors.New("compressor: one of -c or -d is required"))
	}
	args := flag.Args()
	if len(args) == 0 {
		panic(errors.New("compressor: input file is required"))
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		panic(err)
	}
	if *cPtr {
		compressed := cmp.Compress(string(content))
		err = os.WriteFile(*outPtr, []byte(compressed), 0777)
		if err != nil {
			panic(err)
		}
	} else {
		expanded := cmp.Expand(string(content))
		err = os.WriteFile(*outPtr, []byte(expanded), 0777)
		if err != nil {
			panic(err)
		}
	}

}
