package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: program <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	countFrequency(scanner)
}

func countFrequency(scanner *bufio.Scanner) {

	characterMap := make(map[string]int)

	for scanner.Scan() {
		char := rune(scanner.Text()[0])
		counter(char, characterMap)
	}

	for k, v := range characterMap {
		fmt.Printf("%s, occurs %d times \n", k, v)
	}
}

// Increases the value of a map on occurence of a key
func counter(char rune, charMap map[string]int) {
	charMap[strconv.QuoteRuneToASCII(char)]++
}
