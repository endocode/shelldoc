package run

import (
	"fmt"
	"log"
	"math"

	"github.com/endocode/shelldoc/pkg/junitxml"
	"github.com/endocode/shelldoc/pkg/shell"
	"github.com/endocode/shelldoc/pkg/tokenizer"
	"github.com/endocode/shelldoc/pkg/version"
)

func max(a, b int) int { // really, golang?
	if a > b {
		return a
	}
	return b
}

const (
	returnSuccess = iota // the test succeeded
	returnFailure        // the test failed (a problemn with the test)
	returnError          // there was an error executing the test (a problem with shelldoc)
)

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

func (context *Context) performInteractions(inputfile string) (junitxml.JUnitTestSuite, error) {
	// detect shell
	shellpath, err := shell.DetectShell(context.ShellName)
	if err != nil {
		return junitxml.JUnitTestSuite{}, err
	}
	// start a background shell, it will run until the function ends
	shell, err := shell.StartShell(shellpath)
	if err != nil {
		return junitxml.JUnitTestSuite{}, fmt.Errorf("unable to start shell: %v", err)
	}
	defer shell.Exit()
	// read input data
	data, err := ReadInput([]string{inputfile})
	if err != nil {
		return junitxml.JUnitTestSuite{}, fmt.Errorf("unable to read input data: %v", err)
	}
	// run the input through the tokenizer
	visitor := tokenizer.NewInteractionVisitor()
	tokenizer.Tokenize(data, visitor)
	// the test suite object for this file
	suite := junitxml.JUnitTestSuite{}
	suite.Name = inputfile
	suite.AddProperty("shelldoc-version", version.Version())
	// execute the interactions and verify the results:
	fmt.Printf("SHELLDOC: doc-testing \"%s\" ...\n", inputfile)
	// construct the opener and closer format strings, since they depend on verbose mode
	magnitude := int(math.Log10(float64(len(visitor.Interactions)))) + 1
	openerLineEnding := "  : "
	resultString := " "
	if context.Verbose {
		openerLineEnding = "\n"
		resultString = " <-- "
	}
	counterFormat := fmt.Sprintf("%%%ds", magnitude+2)
	opener := fmt.Sprintf(" CMD %s: %%s%s", counterFormat, openerLineEnding)
	closer := fmt.Sprintf("%s%%s\n", resultString)

	for index, interaction := range visitor.Interactions {
		testcase := junitxml.JUnitTestCase{
			Name:      interaction.Cmd,
			Classname: inputfile,
		}
		fmt.Printf(opener, fmt.Sprintf("(%d)", index+1), interaction.Describe())
		if context.Verbose {
			fmt.Printf(" --> %s\n", interaction.Cmd)
		}
		if err := interaction.Execute(&shell); err != nil {
			fmt.Printf(" --  ERROR: %v", err)
			context.RegisterReturnCode(returnError)
			testcase.RegisterFailure(result(returnError), interaction.Result())
		}
		fmt.Printf(closer, interaction.Result())
		if interaction.HasFailure() {
			context.RegisterReturnCode(returnFailure)
			testcase.RegisterFailure(result(returnFailure), interaction.Result())
		}
		suite.RegisterTestCase(testcase)
		if interaction.HasFailure() && context.FailureStops {
			log.Printf("Stop requested after first failed test.")
			break
		}
	}
	fmt.Printf("%s: %d tests - %d successful, %d failures (%d execution errors)\n", result(context.ReturnCode()), suite.TestCount(),
		suite.SuccessCount(), suite.FailureCount(), suite.FailureCountForType(result(returnError)))
	context.Suites.Suites = append(context.Suites.Suites, suite)
	return suite, nil
}
