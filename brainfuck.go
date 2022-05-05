package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/term"
)

// This function will allow us to map all of
// our loop start and stop positions.
func get_loops(input string) map[int]int {
	// Essentially a queue (FILO) to place and pop loop start positions.
	loop_tracker := make(chan int, 100)
	// Map to hold loop start and stop positions.
	loops := map[int]int{}

	// Iterate over the inputted code looking for loops.
	for i, r := range input {
		// Place the loop start position in the queue.
		if r == '[' {
			loop_tracker <- i
		}
		// Grab the latest start position. Map the loop start and stop
		// positions to eachother.
		if r == ']' {
			start := <-loop_tracker
			end := i
			loops[start] = end
			loops[end] = start
		}
	}
	return loops
}

// Driver function that will handle our BF code.
func run(instructions string) {
	// Get all of the loops from the instructions.
	loops := get_loops(instructions)
	// Array of registers to store our data. Starts
	// at 30,000 registers, but will be expanded as
	// necessary.
	registers := []int{}
	for i := 0; i < 30000; i++ {
		registers = append(registers, 0)
	}
	// Pointers for both our instructions and registers.
	inst_ptr := 0
	reg_ptr := 0

	for inst_ptr < len(instructions) {
		switch instructions[inst_ptr] {
		// There is no default case or code cleaning/verification.
		// If an invalid instruction is found it is skipped and ignored.
		case '>':
			// Move to the next register.
			reg_ptr++
			// If we are out of registers, add more.
			if reg_ptr == len(registers) {
				registers = append(registers, 0)
			}
		case '<':
			// Move to the previous register.
			reg_ptr--
			// No negative registers.
			if reg_ptr < 0 {
				reg_ptr = 0
			}
		case '+':
			// Increment register value.
			registers[reg_ptr]++
		case '-':
			// Decrement the register value.
			registers[reg_ptr]--
		case '.':
			// Print the register value to the console. ASCII val -> string/rune
			fmt.Printf("%s", string(registers[reg_ptr]))
		case ',':
			// Read a single byte from user input.
			original_state, err := term.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			defer term.Restore(int(os.Stdin.Fd()), original_state)
			b := make([]byte, 1)
			_, err = os.Stdin.Read(b)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			registers[reg_ptr] = int(b[0])
		case '[':
			// If we are at the start of a loop and the value of the
			// current register is 0, move to after the loop.
			if registers[reg_ptr] == 0 {
				inst_ptr = loops[inst_ptr]
			}
		case ']':
			// If we are at the end of a loop and the value of the
			// current register isn't 0, move to the front of the loop.
			if registers[reg_ptr] != 0 {
				inst_ptr = loops[inst_ptr]
			}
		}
		// Move to the next instruction.
		inst_ptr++
	}
	fmt.Println()
}

func main() {
	// Get input file from user and ensure it is of the correct format.
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./brainfuck <filename.b>")
		os.Exit(0)
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".b") {
		fmt.Println("Error: Improper filename. Ensure that inputted file is of the type <filename.b>")
	}

	// Read inputted file contents.
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error: Could not read file - %s\n", filename)
	}

	// Run the code from the inputted file.
	run(string(contents))
}