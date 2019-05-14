package run

import "github.com/endocode/shelldoc/pkg/junitxml"

// Context contains the context of an exewcution of the run subcommand.
type Context struct {
	// input (configuration) variables
	ShellName     string
	Verbose       bool
	FailureStops  bool
	XMLOutputFile string
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
