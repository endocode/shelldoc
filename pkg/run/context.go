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
	Suites junitxml.JUnitTestSuites
}
