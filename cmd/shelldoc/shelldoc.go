package main

import (
	"io/ioutil"
	golog "log"
	"os"

	"github.com/Endocode/shelldoc/pkg/shell"
	"github.com/Endocode/shelldoc/pkg/tokenizer"
	"github.com/spf13/pflag"
)

const (
	returnSuccess = iota
	returnFailure
	returnError
)

// Options contains the context of a program invocation
type Options struct {
	shell   string // The shell to invoke
	verbose bool   // Enable trace log output
}

// global variables
var (
	options Options
	// Debug receives log messages in verbose mode
	Debug *golog.Logger
	// Log is the standard logger
	Log *golog.Logger
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func result(code int) string {
	switch code {
	case returnFailure:
		return "FAILURE"
	case returnError:
		return "ERROR"
	default:
		return "SUCCESS"
	}
}

func main() {
	pflag.StringVarP(&options.shell, "shell", "s", "", "The shell to invoke (default: the user's shell).")
	pflag.BoolVarP(&options.verbose, "verbose", "v", false, "Enable diagnostic log output.")
	pflag.Parse()
	if options.verbose {
		Debug = golog.New(os.Stderr, "", 0)
	} else {
		Debug = golog.New(ioutil.Discard, "", 0)
	}
	Log = golog.New(os.Stderr, "", 0)
	args := pflag.Args()

	data, err := ReadInput(args)
	if err != nil {
		Log.Fatalln(err)
	}
	// run the input through the tokenizer
	visitor := tokenizer.NewInteractionVisitor()
	tokenizer.Tokenize(data, visitor)
	// start a background shell, it will run until the program ends
	shell, err := shell.StartShell()
	if err != nil {
		Log.Fatalln(err)
	}
	defer shell.Exit()
	returncode := returnSuccess
	testCount := 0
	successCount := 0
	failureCount := 0
	errorCount := 0
	// execute the interactions and verify the results:
	for index, interaction := range visitor.Interactions {
		testCount++
		Log.Printf("COMMAND (%d): %s\n", index+1, interaction.Describe())
		Debug.Printf("--> %s\n", interaction.Cmd)
		if err := interaction.Execute(&shell); err != nil {
			Log.Printf("--  ERROR: %v\n", err)
			returncode = max(returncode, returnError)
			errorCount++
		}
		Debug.Printf("<-- %s\n", interaction.Result())
		if interaction.HasFailure() {
			returncode = max(returncode, returnFailure)
			failureCount++
		} else {
			successCount++
		}
	}
	Log.Printf("%s: %d tests (%d successful, %d failures, %d execution errors)\n", result(returncode), testCount, successCount, failureCount, errorCount)
	os.Exit(returncode)
}
