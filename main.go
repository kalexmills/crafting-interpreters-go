package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		fmt.Errorf("Usage: clox [path]\n")
		os.Exit(64)
	}
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println()
			os.Exit(0)
		}
		if err != nil {
			fmt.Printf("could not read from stdin: %v", err)
			os.Exit(1)
		}
		fmt.Println()
		interpret(line)
	}
}

func runFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("could not read file %s: %v", path, err)
		os.Exit(1)
	}
	result := interpret(string(source))
	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}
}
