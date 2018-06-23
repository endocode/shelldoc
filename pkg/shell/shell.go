package shell

// This file is part of shelldoc.
// © 2018, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: LGPL-3.0

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Shell represents the shell process that runs in the background and executes the commands.
type Shell struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

// DetectShell returns the path to the selected shell or the content of $SHELL
func DetectShell(selected string) (string, error) {
	if len(selected) > 0 {
		// accept what the user said
		log.Printf("Using user-specified shell %s.", selected)
	} else {
		selected = os.Getenv("SHELL")
		log.Printf("Using shell %s (according to $SHELL).", selected)
	}
	if _, err := os.Stat(selected); os.IsNotExist(err) {
		return "", fmt.Errorf("the selected shell does not exist: %v", err)
	}
	return selected, nil
}

// StartShell starts a shell as a background process
func StartShell(shell string) (Shell, error) {
	cmd := exec.Command(shell)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return Shell{}, fmt.Errorf("Unable to set up input stream for shell %s: %v", shell, err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Shell{}, fmt.Errorf("Unable to set up output stream for shell %s: %v", shell, err)
	}
	err = cmd.Start()
	if err != nil {
		return Shell{}, fmt.Errorf("Unable to start shell %s: %v", shell, err)
	}
	return Shell{cmd, stdin, stdout}, nil
}

// ExecuteCommand runs a command in the shell and returns its output and exit code
func (shell *Shell) ExecuteCommand(command string) ([]string, int, error) {
	const (
		beginMarker = ">>>>>>>>>>SHELLDOC_MARKER>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
		endMarker   = "<<<<<<<<<<SHELLDOC_MARKER"
	)
	instruction := fmt.Sprintf("%s\n", strings.TrimSpace(command))
	io.WriteString(shell.stdin, fmt.Sprintf("echo \"%s\"\n", beginMarker))
	io.WriteString(shell.stdin, instruction)
	io.WriteString(shell.stdin, fmt.Sprintf("echo \"%s $?\"\n", endMarker))

	// read output (TODO: with timeout), watch for markers:
	beginEx := fmt.Sprintf("^%s$", beginMarker)
	beginRx := regexp.MustCompile(beginEx)
	endEx := fmt.Sprintf("^%s (.+)$", endMarker)
	endRx := regexp.MustCompile(endEx)

	var output []string
	var rc int
	beginFound := false
	scanner := bufio.NewScanner(shell.stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if beginRx.MatchString(line) {
			beginFound = true
			continue
		}
		if beginFound == false {
			continue
		}
		match := endRx.FindStringSubmatch(line)
		if len(match) > 1 {
			value, err := strconv.Atoi(match[1])
			if err != nil {
				return nil, -1, fmt.Errorf("unable to read exit code for shell command: %v", err)
			}
			rc = value
			break
		}
		output = append(output, line)
	}
	return output, rc, nil
}

// Exit tells a running shell to exit and waits for it
func (shell *Shell) Exit() error {
	io.WriteString(shell.stdin, "exit\n")
	return shell.cmd.Wait()
}
