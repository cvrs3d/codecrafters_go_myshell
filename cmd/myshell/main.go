package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"errors"
)

// type command
func isBuiltin(command string) bool {
    builtins := map[string]bool{
        "echo": true,
        "type": true,
        "exit": true,
    }
    return builtins[command]
}

// Parse type <command>
func handleTypeCmd(input string) error {
    parts := strings.SplitN(input, " ", 2)
    if len(parts) != 2 {
        return errors.New("type: missing operand")
    }

    command := parts[1]

    if isBuiltin(command) {
        fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", command)
    } else {
        fmt.Fprintf(os.Stdout, "%s: not found\n", command)
    }
    return nil
}

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
        }

        if strings.HasPrefix(usrInput, "type ") {
            err := handleTypeCmd(usrInput)
            if err == nil {
                continue
            }
        }

        if strings.HasPrefix(usrInput, "echo ") {
            fmt.Fprintf(os.Stdout,"%s\n", usrInput[5:])
            continue
        }

        fmt.Fprintf(os.Stdout, "%s: not found\n", usrInput)
	}
}
