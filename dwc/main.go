package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	bytesFlag bool
	linesFlag bool
	wordsFlag bool
	charsFlag bool
)

func main() {
	flag.BoolVar(&bytesFlag, "c", false, "number of bytes in a file")
	flag.BoolVar(&linesFlag, "l", false, "number of lines in a file")
	flag.BoolVar(&wordsFlag, "w", false, "number of words in a file")
	flag.BoolVar(&charsFlag, "m", false, "number of characters in a file")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		// No file specified, read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		processInput(scanner)
	} else {
		// File specified
		for _, filePath := range args {
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "File %s does not exist\n", filePath)
				continue
			}
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
				continue
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			processInput(scanner)
			fmt.Printf(" %s\n", filePath)
		}
	}
}

func processInput(scanner *bufio.Scanner) {
	var lines, words, chars, bytes int

	for scanner.Scan() {
		line := scanner.Text()
		lines++
		words += len(splitWords(line))
		chars += len(line)
		bytes += len([]byte(line))
	}

	if linesFlag {
		fmt.Printf("%8d", lines)
	}
	if wordsFlag {

		fmt.Printf("%8d", words)
	}
	if charsFlag {
		fmt.Printf("%8d", chars)
	}
	if bytesFlag {
		fmt.Printf("%8d", bytes)
	}
}

func splitWords(input string) []string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}
