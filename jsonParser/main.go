package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: program <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	if isValidJSON(scanner) {
		fmt.Println("The file contains valid JSON.")
	} else {
		fmt.Println("The file does not contain valid JSON.")
		os.Exit(1)
	}
}

func isValidJSON(scanner *bufio.Scanner) bool {
	for scanner.Scan() {
		char := rune(scanner.Text()[0])
		if !unicode.IsSpace(char) {
			return parseValue(scanner, char)
		}
	}
	return false
}

func parseValue(scanner *bufio.Scanner, firstChar rune) bool {
	switch firstChar {
	case '{':
		return parseObject(scanner)
	case '[':
		return parseArray(scanner)
	case '"':
		return parseString(scanner)
	case 't', 'f', 'n':
		return parseLiteral(scanner, firstChar)
	default:
		if unicode.IsDigit(firstChar) || firstChar == '-' {
			return parseNumber(scanner, firstChar)
		}
	}
	return false
}

func parseObject(scanner *bufio.Scanner) bool {
	// Implementation for parsing objects
	expectKey := true
	for scanner.Scan() {
		char := rune(scanner.Text()[0])
		if unicode.IsSpace(char) {
			continue
		}
		if char == '}' {
			return !expectKey
		}
		if expectKey {
			if char != '"' {
				return false
			}
			if !parseString(scanner) {
				return false
			}
			if !expectChar(scanner, char) {
				return false
			}
			expectKey = false
		} else {
			if !parseValue(scanner, char) {
				return false
			}
			if !expectChar(scanner, ',', '}') {
				return false
			}
			expectKey = (scanner.Text()[0] == ',')
		}
	}
	return false // Placeholder
}

func parseArray(scanner *bufio.Scanner) bool {
	// Implementation for parsing arrays
	return true // Placeholder
}

func parseString(scanner *bufio.Scanner) bool {
	// Implementation for parsing strings
	for scanner.Scan() {
		char := rune(scanner.Text()[0])
		if char == '"' {
			return true
		}
	}
	return false // Placeholder
}

func parseLiteral(scanner *bufio.Scanner, firstChar rune) bool {
	// Implementation for parsing true, false, null
	return true // Placeholder
}

func parseNumber(scanner *bufio.Scanner, firstChar rune) bool {
	// Implementation for parsing numbers
	return true // Placeholder
}

func expectChar(scanner *bufio.Scanner, expected ...rune) bool {
	if !scanner.Scan() {
		return false
	}
	char := rune(scanner.Text()[0])
	for _, e := range expected {
		if char == e {
			return true
		}
	}
	return false
}
