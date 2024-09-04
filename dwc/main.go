package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	ytes  string
	lines string
	words string
)

func main() {
	flag.StringVar(&ytes, "c", "", "number of bytes in a file")
	flag.StringVar(&lines, "l", "", "number of lines in a file")
	flag.StringVar(&words, "w", "", "number of words in a file")
	flag.Parse()

	if ytes != "" {
		ByteCount(ytes)
	} else if lines != "" {
		LineCount(lines)
	} else if words != "" {
		WordsCount(words)
	}
}

func checkFile(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func ByteCount(input string) {
	isFile := checkFile(input)
	if isFile {
		content, err := os.ReadFile(input)
		if err != nil {
			fmt.Printf("%d", len(input))
		}
		fmt.Printf("%d %s", len(content), input)
	} else {
		fmt.Printf("%d %s", len(input), input)
	}

}

func LineCount(input string) {
	isFile := checkFile(input)
	if isFile {
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
	} else {
		fmt.Println("Hello")
	}
}

func WordsCount(input string) {
	isFile := checkFile(input)
	if isFile {
		content, err := os.Open(input)
		if err != nil {
			fmt.Println(err)
		}
		defer content.Close()

		fileScanner := bufio.NewScanner(content)
		fileScanner.Split(bufio.ScanWords)
		var fileWords []string

		for fileScanner.Scan() {
			fileWords = append(fileWords, fileScanner.Text())
		}

		fmt.Printf("%d %s", len(fileWords), input)
	}
}
