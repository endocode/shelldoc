package tokenizer

// This file is part of shelldoc.
// © 2018, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: LGPL-3.0

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/endocode/shelldoc/pkg/shell"
)

const (
	// NewInteraction indicates that the interaction has not been executed yet
	NewInteraction = iota
	// ResultExecutionError indicates that there has been an error in executing the command, not with the command itself
	ResultExecutionError
	// ResultError indicates that the command exited with an non-zero exit code
	ResultError
	// ResultMatch means the output directly matched the expected output
	ResultMatch
	// ResultRegexMatch means the output matched the alternative regex
	ResultRegexMatch
	// ResultMismatch indicates that the output from the command did not match expectations in any way
	ResultMismatch
)

// Interaction represents one interaction with the shell
type Interaction struct {
	// Cmd contains exactly the command the shell is supposed to execute
	Cmd string
	// Response contains the exected response from the shell, in plain text
	Response []string
	//AlternativeRegEx string
	// Language contains the language specified if the interaction was extracted from a fenced code block
	Language string
	// Attributes contains the shelldoc attributes specified in a fenced code block
	Attributes map[string]string
	// Caption contains a descriptive name for the interaction
	Caption string
	// Result contains a human readable description of the result after the interaction has been executed
	ResultCode int
	// Comment contains an explanation of the ResultCode after execution
	Comment string
}

// Describe returns a human-readable description of the interaction
func (interaction *Interaction) Describe() string {
	const elideCmdAt = 40
	const elideResponseAt = 25
	format := fmt.Sprintf("%%-%ds  ?  %%-%ds", elideCmdAt, elideResponseAt)
	name := interaction.Cmd
	if len(interaction.Caption) != 0 {
		name = interaction.Caption
	}
	expect := elideString(strings.Join(interaction.Response, ", "), elideResponseAt)
	if len(expect) == 0 {
		expect = "(no response expected)"
	}
	result := fmt.Sprintf(format, elideString(name, elideCmdAt), expect)
	return result
}

// Result returns a human readable description of the result of the interaction
func (interaction *Interaction) Result() string {
	switch interaction.ResultCode {
	case NewInteraction:
		return "not executed"
	case ResultExecutionError:
		return "ERROR (result not evaluated)"
	case ResultMatch:
		if len(interaction.Response) == 0 {
			return "PASS (execution successful)"
		}
		return "PASS (match)"
	case ResultRegexMatch:
		return "PASS (regex match)"
	case ResultMismatch:
		return "FAIL (mismatch)"
	case ResultError:
		return "FAIL (execution failed)"
	default:
		return "YOU FOUND A BUG!!11!1!"
	}
}

// HasFailure returns true if the interaction failed (not on execution errors)
func (interaction *Interaction) HasFailure() bool {
	return interaction.ResultCode == ResultError || interaction.ResultCode == ResultMismatch
}

// New creates an empty interaction with a Caption
func New(caption string) *Interaction {
	interaction := new(Interaction)
	interaction.Caption = caption
	return interaction
}

// evaluateResponse compares the output to the expected response, and respects "ellipsis" (don't care from here on forward)
func (interaction *Interaction) evaluateResponse(response []string) bool {
	output := response
	expected := interaction.Response
	for index, line := range interaction.Response {
		if strings.TrimSpace(line) == "..." {
			output = response[:index]
			expected = interaction.Response[:index]
			break
		}
	}
	if len(output) == 0 && len(expected) == 0 {
		return true
	}
	return reflect.DeepEqual(output, expected)
}

// Execute the interaction and store the result
func (interaction *Interaction) Execute(shell *shell.Shell) error {
	// execute the command in the shell
	output, rc, err := shell.ExecuteCommand(interaction.Cmd)
	// compare the results
	if err != nil {
		interaction.ResultCode = ResultExecutionError
		interaction.Comment = err.Error()
		return fmt.Errorf("unable to execute command: %v", err)
	} else if rc != 0 {
		interaction.ResultCode = ResultError
		interaction.Comment = fmt.Sprintf("command exited with non-zero exit code %d", rc)
	} else if interaction.evaluateResponse(output) {
		interaction.ResultCode = ResultMatch
		interaction.Comment = ""
	} else if interaction.compareRegex(output) {
		interaction.ResultCode = ResultRegexMatch
	} else {
		interaction.ResultCode = ResultMismatch
		interaction.Comment = ""
	}
	return nil
}

func (interaction *Interaction) compareRegex(output []string) bool {
	// match, err := regexp.MatchString(interaction.AlternativeRegEx, output); err
	return false
}

func elideString(text string, length int) string {
	if length > 6 && len(text) > length {
		return fmt.Sprintf("%s...", text[:length-3])
	}
	return text
}
