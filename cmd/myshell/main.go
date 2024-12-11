package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	reader := bufio.NewReader(os.Stdin)

	usrInput, err := reader.ReadString('\n')
	if err != nil {
	    fmt.Fprint(os.Stdout, "invalid_command: not found\n")
	}
    usrInput = usrInput[:len(usrInput) - 1]
	fmt.Fprintf(os.Stdout, "%s: not found\n", usrInput)
}
