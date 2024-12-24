package main

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// parseInput handles input string and works with ''
func parseInput(input string) ([]string, map[string]string, error) {
    args := []string{}
    redirects := make(map[string]string)
    var currentArg strings.Builder
    inSingleQuote, inDoubleQuote, escapeNext := false, false, false

    for i := 0; i < len(input); i++ {
        char := input[i]

        if escapeNext {
            currentArg.WriteByte(char)
            escapeNext = false
            continue
        }

        switch char {
        case '\\':
            escapeNext = true
        case '\'':
            if !inDoubleQuote {
                inSingleQuote = !inSingleQuote
            } else {
                currentArg.WriteByte(char)
            }
        case '"':
            if !inSingleQuote {
                inDoubleQuote = !inDoubleQuote
            } else {
                currentArg.WriteByte(char)
            }
        case '>':
            if !inSingleQuote && !inDoubleQuote {
                if currentArg.Len() > 0 {
                    args = append(args, currentArg.String())
                    currentArg.Reset()
                }
                targetFile := strings.TrimSpace(input[i+1:])
                if len(targetFile) == 0 {
                    return nil, nil, errors.New("parse error: no file specified for redirection")
                }
                if i > 0 && input[i-1] == '1' {
                    redirects["stdout"] = targetFile
                    i++ // Skip the '1' character
                } else {
                    redirects["stdout"] = targetFile
                }
            } else {
                currentArg.WriteByte(char)
            }
        case ' ':
            if !inSingleQuote && !inDoubleQuote {
                if currentArg.Len() > 0 {
                    args = append(args, currentArg.String())
                    currentArg.Reset()
                }
            } else {
                currentArg.WriteByte(char)
            }
        default:
            currentArg.WriteByte(char)
        }
    }

    if inSingleQuote || inDoubleQuote {
        return nil, nil, errors.New("parse error: mismatched quotes")
    }

    if currentArg.Len() > 0 {
        args = append(args, currentArg.String())
    }

    return args, redirects, nil
}

// Checks if command is a builtin
func isBuiltin(command string) bool {
    builtins := map[string]bool{
        "echo": true,
        "type": true,
        "exit": true,
        "pwd":  true,
        "cd":   true,
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

        if fileInfo, err := os.Stat(fullPath); err == nil && !fileInfo.IsDir() {
            return fullPath, nil
        }
    }

    return "", errors.New(fmt.Sprintf("%s: not found\n", command))
}

// Parse type <command>
func handleTypeCmd(args []string) error {
    if len(args) != 1 {
        return errors.New("type: missing operand")
    }

    command := args[0]

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
func executeCommand(args []string, redirects map[string]string) error {
    if len(args) == 0 {
        return errors.New("no command provided")
    }

    command := args[0]
    cmd := exec.Command(command, args[1:]...)

    // Handle stdout redirection
    if stdoutFile, ok := redirects["stdout"]; ok {
        file, err := os.OpenFile(stdoutFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
            return fmt.Errorf("failed to open file for stdout redirection: %v", err)
        }
        defer file.Close()
        cmd.Stdout = file
    } else {
        cmd.Stdout = os.Stdout
    }

    // Handle stderr redirection
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("command execution failed: %v", err)
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
        return fmt.Errorf("cd: %v", err)
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
            return fmt.Errorf("cat: %v", err)
        }

        content, err := io.ReadAll(file)
        if err != nil {
            file.Close()
            return fmt.Errorf("cat: %v", err)
        }
        file.Close()

        result.Write(content)
    }
    fmt.Fprint(os.Stdout, result.String())
    return nil
}

func main() {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Fprint(os.Stdout, "$ ")

        usrInput, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprint(os.Stdout, "invalid_command: not found\n")
            continue
        }

        usrInput = strings.TrimSpace(usrInput)

        if usrInput == "exit 0" {
            break
        }

        args, redirects, err := parseInput(usrInput)
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
            err = handleTypeCmd(args)
        case "echo":
            handleEcho(args)
        case "cat":
            err = handleCat(args)
        case "pwd":
            err = handlePwd()
        case "cd":
            err = handleCd(usrInput)
        default:
            err = executeCommand(append([]string{command}, args...), redirects)
        }

        if err != nil {
            fmt.Println(err)
        }
    }
}