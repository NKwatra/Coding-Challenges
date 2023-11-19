package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("ncut: %v", err))
	}
}

func main() {
	defer func() {
		text := recover()
		if val, ok := text.(string); ok {
			os.Stdout.WriteString(val + "\n")
			os.Exit(1)
		}
	}()
	var delimeter string
	fieldPtr := flag.String("f", "", "List of fields to return")
	sPtr := flag.Bool("s", false, "To suppress lines with no delimeter")
	flag.StringVar(&delimeter, "d", "\t", "-d ,")
	flag.Parse()
	var file *os.File
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(flag.Arg(0))
		checkErr(err)
	}
	printFields(file, delimeter, *fieldPtr, *sPtr)
}

func printFields(file *os.File, delim, field string, s bool) {
	var splitDelim string
	if strings.Contains(field, " ") {
		splitDelim = " "
	} else {
		splitDelim = ","
	}
	fieldSegments := strings.Split(field, splitDelim)
	columns := make([]int, len(fieldSegments))
	var err error
	for i, fieldSegment := range fieldSegments {
		columns[i], err = strconv.Atoi(fieldSegment)
		if err != nil {
			checkErr(err)
		}
	}
	reader := bufio.NewReader(file)
	var line []byte
	fieldNumber := slices.Max(columns)
	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			break
		}
		segments := strings.Split(string(line), delim)
		if len(segments) < fieldNumber {
			if !s {
				os.Stdout.Write(line)
			}
		} else {
			values := make([]string, len(columns))
			for i, index := range columns {
				values[i] = segments[index-1]
			}
			os.Stdout.WriteString(strings.Join(values, delim))
		}
		os.Stdout.WriteString("\n")
	}
}
