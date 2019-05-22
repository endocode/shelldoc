package run

import (
	"fmt"
	"os"

	"github.com/endocode/shelldoc/pkg/junitxml"
)

// Context contains the context of an exewcution of the run subcommand.
type Context struct {
	// input (configuration) variables
	ShellName     string
	Verbose       bool
	FailureStops  bool
	XMLOutputFile string
	ReplaceDots   bool
	Files         []string
	// output variables
	Suites     junitxml.JUnitTestSuites
	returnCode int
}

// RegisterReturnCode registers a potential error. The return code can never decrease.
func (context *Context) RegisterReturnCode(returnCode int) int {
	context.returnCode = max(context.returnCode, returnCode)
	return context.returnCode
}

// ReturnCode returns the overall result of the operation.
func (context *Context) ReturnCode() int {
	return context.returnCode
}

// WriteXML writes the test results to the specified XML output file
func (context *Context) WriteXML() error {
	if len(context.XMLOutputFile) > 0 {
		file, err := os.OpenFile(context.XMLOutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return fmt.Errorf("unable to open XML output file for writing: %v", err)
		}
		if err := context.Suites.Write(file); err != nil {
			return fmt.Errorf("error writing XML output file: %v", err)
		}
	}
	return nil
}

// ExecuteFiles runs each file through performInteractions and aggregates the results
func (context *Context) ExecuteFiles() int {
	context.RegisterReturnCode(returnSuccess)
	for _, file := range context.Files {
		suite, err := context.performInteractions(file)
		if err != nil {
			fmt.Println(err) // log may be disabled (see "verbose")
			os.Exit(returnError)
		}
		context.Suites.Suites = append(context.Suites.Suites, *suite)
	}
	if err := context.WriteXML(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(returnError)
	}
	return context.ReturnCode()
}
