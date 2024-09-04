package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	words string
	lines string
)

func main() {
	flag.StringVar(&words, "c", "", "number of bytes in a file")
	flag.StringVar(&lines, "l", "", "number of lines in a file")
	flag.Parse()

	if words != "" {
		ByteCount(words)
	}
	if lines != "" {
		LineCount(lines)
	}
}

func ByteCount(input string) {
	content, err := os.ReadFile(input)
	if err != nil {
		fmt.Printf("%d", len(input))
	}
	fmt.Printf("%d %s", len(content), input)
}

func LineCount(input string) {
	content, err := os.Open(input)
	if err != nil {
		fmt.Println(err)
	}
	defer content.Close()

	fileScanner := bufio.NewScanner(content)
	fileScanner.Split(bufio.ScanLines)

	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	fmt.Printf("%d %s", len(fileLines), input)
}
