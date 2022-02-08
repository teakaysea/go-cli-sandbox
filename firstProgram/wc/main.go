package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	countWords = iota
	countLines
	countBytes
)

func main() {
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	flag.Parse()
	countType := countWords
	if *lines {
		countType = countLines
	} else if *bytes {
		countType = countBytes
	}
	fmt.Println(count(os.Stdin, countType))
}

func count(r io.Reader, countType int) int {
	scanner := bufio.NewScanner(r)
	switch countType {
	case countWords:
		scanner.Split(bufio.ScanWords)
	case countBytes:
		scanner.Split(bufio.ScanBytes)
	}

	var wc int
	for scanner.Scan() {
		wc++
	}
	return wc
}
