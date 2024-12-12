package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {


	// Wait for user input
	reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Fprint(os.Stdout, "$ ")
        usrInput, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprint(os.Stdout, "invalid_command: not found\n")
            continue
        }
        usrInput = strings.TrimSpace(usrInput)
        if usrInput == "exit 0"{
            break
        fmt.Fprintf(os.Stdout, "%s: not found\n", usrInput)
	}
}
