package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/endocode/shelldoc/pkg/shell"
	"github.com/endocode/shelldoc/pkg/tokenizer"
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
var options Options

func max(a, b int) int { // really, golang?
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
	// verbose essentially enables or disables log output:
	if options.verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(0)
	log.SetPrefix("Note: ")
}

func performInteractions(inputfile string) (resultStats, error) {
	// detect shell
	shellpath, err := shell.DetectShell(options.shell)
	if err != nil {
		return resultStats{}, err
	}

	// start a background shell, it will run until the function ends
	shell, err := shell.StartShell(shellpath)
	if err != nil {
		return resultStats{}, fmt.Errorf("unable to start shell: %v", err)
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
	// construct the opener and closer format strings, since they depend on verbose mode
	magnitude := int(math.Log10(float64(len(visitor.Interactions)))) + 1
	openerLineEnding := "  : "
	resultString := " "
	if options.verbose {
		openerLineEnding = "\n"
		resultString = " <-- "
	}
	counterFormat := fmt.Sprintf("%%%ds", magnitude+2)
	opener := fmt.Sprintf(" CMD %s: %%s%s", counterFormat, openerLineEnding)
	closer := fmt.Sprintf("%s%%s\n", resultString)

	for index, interaction := range visitor.Interactions {
		results.testCount++
		fmt.Printf(opener, fmt.Sprintf("(%d)", index+1), interaction.Describe())

		if options.verbose {
			fmt.Printf(" --> %s\n", interaction.Cmd)
		}
		if err := interaction.Execute(&shell); err != nil {
			fmt.Printf(" --  ERROR: %v", err)
			results.returncode = max(results.returncode, returnError)
			results.errorCount++
		}
		fmt.Printf(closer, interaction.Result())
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
	pflag.StringVarP(&options.shell, "shell", "s", "", "The shell to invoke (default: $SHELL).")
	pflag.BoolVarP(&options.verbose, "verbose", "v", false, "Enable diagnostic log output.")
	pflag.Parse()
	initializeLogging()
	args := pflag.Args()
	returnCode := returnSuccess
	for _, file := range args {
		results, err := performInteractions(file)
		if err != nil {
			fmt.Println(err) // log may be disabled (see "verbose")
			os.Exit(returnError)
		}
		returnCode = max(results.returncode, returnCode)
	}
	os.Exit(returnCode)
}
