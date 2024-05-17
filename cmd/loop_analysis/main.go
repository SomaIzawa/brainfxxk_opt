package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// extractLoops function extracts looped sections from Brainfuck code
func extractLoops(code string) ([]string, error) {
	var loops []string
	var stack []int

	for i, char := range code {
		switch char {
		case '[':
			// Push the index of '[' onto the stack
			stack = append(stack, i)
		case ']':
			// Pop the index of the matching '[' from the stack
			if len(stack) == 0 {
				return nil, fmt.Errorf("unmatched ']' at position %d", i)
			}
			start := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// Append the loop to the result
			loops = append(loops, code[start:i+1])
		}
	}

	// If there are any unmatched '[' left in the stack, return an error
	if len(stack) > 0 {
		return nil, fmt.Errorf("unmatched '[' at position %d", stack[len(stack)-1])
	}

	return loops, nil
}

func main() {
	// Open the file
	file, err := os.Open("./example/mandelbrot.bf")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert content to string and remove any surrounding whitespace
	code := strings.TrimSpace(string(content))

	// Extract loops from the code
	loops, err := extractLoops(code)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the extracted loops
	for i, loop := range loops {
		fmt.Println(i)
		fmt.Println(loop)
	}
}
