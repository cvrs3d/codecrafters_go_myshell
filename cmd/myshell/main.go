package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"errors"
	"path/filepath"
	"io"
)

// parseInput handles input string and works with ''
func parseInput(input string) ([]string, error) {
    var args []string
    var currentArg strings.Builder
    inSingleQuote := false

    for _, char := range input {
        switch {
        case char == '\'':
            inSingleQuote = !inSingleQuote
        case char == ' ' && !inSingleQuote:
            if currentArg.Len() > 0 {
                args = append(args, currentArg.String())
                currentArg.Reset()
            }
        default:
            currentArg.WriteRune(char)
        }
    }

    if inSingleQuote {
        return nil, errors.New("unmatched single quote in input")
    }

    if currentArg.Len() > 0 {
        args = append(args, currentArg.String())
    }

    return args, nil
}

// Checks if command is a builtin
func isBuiltin(command string) bool {
    builtins := map[string]bool{
        "echo": true,
        "type": true,
        "exit": true,
        "pwd": true,
        "cd": true,
        "cat": true,
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

    return "", errors.New(fmt.Sprintf("%s: not found\n", command))
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
        return fmt.Errorf("%s: not found", command)
    }

    // If command is found
    fmt.Fprintf(os.Stdout, "%s is %s\n", command, path)
    return nil
}

// Executes new command
func executeCommand(input string) error {
    args := strings.Fields(input)
    if len(args) == 0 {
        return errors.New("no command provided")
    }

    command := args[0]
    path, err := findCommandInPath(command)
    if err != nil {
        return fmt.Errorf("%s: command not found", command)
    }

    cmd := exec.Command(path, args[1:]...)

    // Redirecting channels
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("%s: command not found", command)
    }

    return nil
}

// Handle the pwd builtin
func handlePwd() error {
    currentDir, err := os.Getwd()

    if err != nil {
        return fmt.Errorf("pwd: %v", err)
    }
    fmt.Println(currentDir)
    return nil
}

// Handle CD
func handleCd(input string) error {
    parts := strings.Fields(input)

    var path string

    if len(parts) == 1 || parts[1] == "~" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return fmt.Errorf("cd: could not get home directory: %v", err)
        }
        path = homeDir
    } else {
        path = parts[1]
    }

    if err := os.Chdir(path); err != nil {
        return fmt.Errorf("cd: %s: No such file or directory", path)
    }

    return nil
}

// handleEcho works with echo
func handleEcho(args []string) {
	if len(args) > 0 {
		// Print joined args space-separated
		fmt.Println(strings.Join(args, " "))
	}
}

// handleCat stands for cat builtin
func handleCat(args []string) error {
    if len(args) == 0 {
        return errors.New("cat: missing file operand")
    }

    var result strings.Builder

    for _, filename := range args {
        file, err := os.Open(filename)
        if err != nil {
            return fmt.Errorf("cat: cannot open '%s': %v", filename, err)
        }

        content, err := io.ReadAll(file)
        if err != nil {
            file.Close()
            return fmt.Errorf("cat: error reading '%s': %v", filename, err)
        }
        file.Close()

        result.Write(content)
    }
    fmt.Fprintln(os.Stdout, result.String())
    return nil
}

func main() {


	// Wait for user input
	reader := bufio.NewReader(os.Stdin)

    for {
        // $$$$$$$
        fmt.Fprint(os.Stdout, "$ ")

        // Reading users input
        usrInput, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprint(os.Stdout, "invalid_command: not found\n")
            continue
        }

        usrInput = strings.TrimSpace(usrInput)

        // If exit
        if usrInput == "exit 0"{
            break
        }

        args, err := parseInput(usrInput)
        if err != nil {
            fmt.Println(err)
            continue
        }

        if len(args) == 0 {
            continue
        }

        command := args[0]
        args = args[1:]

        switch command {
        case "type":
            if err := handleTypeCmd(strings.Join(args, " ")); err != nil {
                fmt.Fprintln(os.Stdout, err)
            }
        case "echo":
            handleEcho(args)
        case "cat":
            if err := handleCat(args); err != nil {
				fmt.Println(err)
			}
        case "pwd":
            if err := handlePwd(); err != nil {
                fmt.Fprintf(os.Stdout, "%s\n", err.Error())
            }
        case "cd":
            if err := handleCd(usrInput); err != nil {
                fmt.Fprintln(os.Stdout, err)
            }
        default:
            if err := executeCommand(usrInput); err != nil {
            fmt.Fprintln(os.Stdout, err)
            }
        }
	}
}
