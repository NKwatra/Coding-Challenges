package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

/*
*
Handle any errors and exit with code 1
*/
func handleError(err error) {
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok {
			fmt.Printf("mywc: %v\n", pathError)
		} else {
			fmt.Println(err)
		}
		panic(1)
	}
}

func main() {
	// recover from an error and exit with integer code
	defer func() {
		code := recover()
		if val, ok := code.(int); ok {
			os.Exit(val)
		}
	}()

	bytePtr := flag.Bool("c", false, "To count the number of bytes in input")
	linesPtr := flag.Bool("l", false, "To count the number of lines in input")
	wordsPtr := flag.Bool("w", false, "To count the number of words in the input")
	charactersPtr := flag.Bool("m", false, "To count the number of characters in the input")
	flag.Parse()
	args := flag.Args()
	var file *os.File
	var err error
	if len(args) == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(args[0])
		handleError(err)
	}
	defer file.Close()
	var count string
	if *bytePtr {
		count = fmt.Sprint(computeBytes(file))
	} else if *linesPtr {
		count = fmt.Sprint(computeLines(file))
	} else if *wordsPtr {
		count = fmt.Sprint(computeWords(file))
	} else if *charactersPtr {
		count = fmt.Sprint(computeCharacters(file))
	} else {
		count = fmt.Sprintf("%v   %v   %v", computeLines(file), computeWords(file), computeCharacters(file))
	}
	if len(args) == 0 {
		fmt.Println(count)
	} else {
		fmt.Printf("    %s %s\n", count, args[0])
	}
}

func computeBytes(file *os.File) int {
	reader := bufio.NewReader(file)
	count := 0
	for {
		line, _, readErr := reader.ReadLine()
		if readErr != nil {
			break
		}
		count += len(line) + 1
	}
	file.Seek(0, 0)
	return count
}

func computeLines(file *os.File) int {
	reader := bufio.NewReader(file)
	count := 0
	for {
		_, _, readErr := reader.ReadLine()
		if readErr != nil {
			break
		}
		count++
	}
	file.Seek(0, 0)
	return count
}

func computeWords(file *os.File) int {
	reader := bufio.NewReader(file)
	count := 0
	for {
		line, _, readErr := reader.ReadLine()
		if readErr != nil {
			break
		}
		count += len(strings.Fields(string(line)))
	}
	file.Seek(0, 0)
	return count
}

func computeCharacters(file *os.File) int {
	reader := bufio.NewReader(file)
	count := 0
	for {
		line, _, readErr := reader.ReadLine()
		if readErr != nil {
			break
		}
		// since readline ignore the newline, append it in
		// counting number of characters
		count += utf8.RuneCountInString(string(line) + "\n")
	}
	file.Seek(0, 0)
	return count
}
