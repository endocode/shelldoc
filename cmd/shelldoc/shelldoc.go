package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	Debug *log.Logger
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

type resultStats struct {
	returncode, testCount, successCount, failureCount, errorCount int
}

func initializeLogging() {
	if options.verbose {
		Debug = log.New(os.Stderr, "", 0)
	} else {
		Debug = log.New(ioutil.Discard, "", 0)
	}
}

func performInteractions(inputfile string) (resultStats, error) {
	// start a background shell, it will run until the function ends
	shell, err := shell.StartShell()
	if err != nil {
		log.Fatalln(err)
	}
	defer shell.Exit()

	// read input data
	data, err := ReadInput([]string{inputfile})
	if err != nil {
		return resultStats{}, fmt.Errorf("unable to read input data: %v", err)
	}
	// run the input through the tokenizer
	visitor := tokenizer.NewInteractionVisitor()
	tokenizer.Tokenize(data, visitor)

	// execute the interactions and verify the results:
	fmt.Printf("SHELLDOC: doc-testing \"%s\" ...\n", inputfile)
	results := resultStats{returnSuccess, 0, 0, 0, 0}
	for index, interaction := range visitor.Interactions {
		results.testCount++
		if options.verbose {
			fmt.Printf(" COMMAND (%d): %s\n", index+1, interaction.Describe())
		} else {
			fmt.Printf(" COMMAND (%d): %s ... ", index+1, interaction.Describe())
		}
		Debug.Printf(" --> %s\n", interaction.Cmd)
		if err := interaction.Execute(&shell); err != nil {
			fmt.Printf(" --  ERROR: %v", err)
			results.returncode = max(results.returncode, returnError)
			results.errorCount++
		}
		if options.verbose {
			fmt.Printf("<-- %s\n", interaction.Result())
		} else {
			fmt.Printf("%s\n", interaction.Result())
		}
		if interaction.HasFailure() {
			results.returncode = max(results.returncode, returnFailure)
			results.failureCount++
		} else {
			results.successCount++
		}
	}
	fmt.Printf("%s: %d tests (%d successful, %d failures, %d execution errors)\n", result(results.returncode), results.testCount, results.successCount, results.failureCount, results.errorCount)
	return results, nil
}

func main() {
	pflag.StringVarP(&options.shell, "shell", "s", "", "The shell to invoke (default: the user's shell).")
	pflag.BoolVarP(&options.verbose, "verbose", "v", false, "Enable diagnostic log output.")
	pflag.Parse()
	initializeLogging()
	args := pflag.Args()
	returnCode := returnSuccess
	for _, file := range args {
		results, err := performInteractions(file)
		if err != nil {
			log.Fatalf("%v", err)
		}
		returnCode = max(results.returncode, returnCode)
	}
	os.Exit(returnCode)
}
