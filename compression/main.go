package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: program <filepath> <outputfilepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	outputFilePath := os.Args[2]

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
	symFreq := countFrequency(scanner)

	tre := buildTrees(symFreq)

	fmt.Println("SYMBOL\tWEIGHT\tHUFFMAN CODE")
	codeTable := printCodes(tre, []byte{})
	for char, code := range codeTable {
		fmt.Printf("%c:\t%s\n", char, code)
	}

	err = encodeText(filePath, outputFilePath, codeTable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding text: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Encoded text has been written to %s\n", outputFilePath)
}

func countFrequency(scanner *bufio.Scanner) map[rune]int {

	characterMap := make(map[rune]int)

	for scanner.Scan() {
		char := rune(scanner.Text()[0])
		counter(char, characterMap)
	}
	return characterMap
}

// Increases the value of a map on occurence of a key
func counter(char rune, charMap map[rune]int) {
	charMap[char]++
}

type HuffNode struct {
	freq        int
	left, right HuffTree
}

type HuffLeaf struct {
	freq    int
	element rune
}

type HuffTree interface {
	Freq() int
}

func (self HuffLeaf) Freq() int {
	return self.freq
}
func (self HuffNode) Freq() int {
	return self.freq
}

type treeHeap []HuffTree

func (th treeHeap) Len() int {
	return len(th)
}

func (th treeHeap) Less(i, j int) bool {
	return th[i].Freq() < th[j].Freq()
}

func (th *treeHeap) Push(elem interface{}) {
	*th = append(*th, elem.(HuffTree))
}

func (th *treeHeap) Pop() (popped interface{}) {
	popped = (*th)[len(*th)-1]
	*th = (*th)[:len(*th)-1]
	return
}
func (th treeHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
}

func buildTrees(symFreq map[rune]int) HuffTree {
	var trees treeHeap
	for c, f := range symFreq {
		trees = append(trees, HuffLeaf{freq: f,
			element: c},
		)
	}
	heap.Init(&trees)

	for trees.Len() > 1 {
		a := heap.Pop(&trees).(HuffTree)
		b := heap.Pop(&trees).(HuffTree)

		heap.Push(&trees, HuffNode{a.Freq() + b.Freq(), a, b})
	}
	return heap.Pop(&trees).(HuffTree)
}

func printCodes(tree HuffTree, prefix []byte) map[rune]string {
	codeTable := make(map[rune]string)
	var traverse func(HuffTree, []byte)
	traverse = func(t HuffTree, prefix []byte) {
		switch i := t.(type) {
		case HuffLeaf:
			fmt.Printf("%c\t%d\t%s\n", i.element, i.freq, string(prefix))
			codeTable[i.element] = string(prefix)
		case HuffNode:
			prefix = append(prefix, '0')
			traverse(i.left, prefix)
			prefix = prefix[:len(prefix)-1]

			prefix = append(prefix, '1')
			traverse(i.right, prefix)
			prefix = prefix[:len(prefix)-1]
		}
	}
	traverse(tree, prefix)
	return codeTable
}

func encodeText(inputFile string, outputFile string, codeTable map[rune]string) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanRunes)

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for scanner.Scan() {
		char := []rune(scanner.Text())[0]
		if code, ok := codeTable[char]; ok {
			_, err := writer.WriteString(code)
			if err != nil {
				return fmt.Errorf("error writing to output file: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	return nil
}
