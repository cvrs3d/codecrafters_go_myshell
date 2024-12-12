package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"errors"
	"path/filepath"
)

// Checks if command is a builtin
func isBuiltin(command string) bool {
    builtins := map[string]bool{
        "echo": true,
        "type": true,
        "exit": true,
    }
    return builtins[command]
}

// Scans the PATH
func findCommandInPath(command string) (string, error) {
    pathEnv, exists := os.LookupEnv("PATH")

    if !exists || pathEnv == "" {
        return "", errors.New("type: PATH env var not found")
    }

    paths := strings.Split(pathEnv, string(os.PathListSeparator))

    // We search for certain command
    for _, dir := range paths {
        fullPath := filepath.Join(dir, command)

        if fileInfo, err := os.Stat(fullPath); err == nil {
            if !fileInfo.IsDir() {
                return fullPath, nil
            }
        }
    }

    return "", errors.New(fmt.Sprintf("%s: not found", command))
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
        return nil
    }

    path, err := findCommandInPath(command)
    if err != nil {
        return err
    }

    // If command is found
    fmt.Fprintf(os.Stdout, "%s is %s\n", command, path)
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
            } else {
                fmt.Fprint(os.Stdout, err)
            }
        }

        if strings.HasPrefix(usrInput, "echo ") {
            fmt.Fprintf(os.Stdout,"%s\n", usrInput[5:])
            continue
        }

        fmt.Fprintf(os.Stdout, "%s: not found\n", usrInput)
	}
}
