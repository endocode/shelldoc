package interaction

import (
	"fmt"

	"github.com/Endocode/shelldoc/pkg/shell"
)

const (
	// NewInteraction indicates that the interaction has not been executed yet
	NewInteraction = iota
	// ResultError indicates that there has been an error in executing the command, not the command itself
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
	// Caption contains a descriptive name for the interaction
	Caption string
	// Result contains a human readable description of the result after the interaction has been executed
	ResultCode int
}

// Describe returns a human-readable description of the interaction
func (interaction *Interaction) Describe() string {
	if len(interaction.Caption) == 0 {
		return fmt.Sprintf("command \"%s\"", interaction.Cmd)
	}
	return interaction.Caption
}

// Result returns a human readable description of the result of the interaction
func (interaction *Interaction) Result() string {
	switch interaction.ResultCode {
	case NewInteraction:
		return "not executed"
	case ResultError:
		return "ERROR (result not evaluated)"
	case ResultMatch:
		return "PASS (match)"
	case ResultRegexMatch:
		return "PASS (regex match)"
	case ResultMismatch:
		return "FAIL (mismatch)"
	}
	return "WTF"
}

// New creates an empty interaction with a Caption
func New(caption string) *Interaction {
	interaction := new(Interaction)
	interaction.Caption = caption
	return interaction
}

// Execute the interaction and store the result
func (interaction *Interaction) Execute(shell *shell.Shell) error {
	//NI
	// ...execute the command in the shell
	// compare the results
	// store
	interaction.ResultCode = ResultError
	return fmt.Errorf("NI")
}
