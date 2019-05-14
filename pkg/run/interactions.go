package run

import (
	"fmt"
	"log"
	"math"
	"os"

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

type resultStats struct {
	successCount, failureCount, errorCount int
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

func (context *Context) performInteractions(inputfile string) (resultStats, error) {
	// detect shell
	shellpath, err := shell.DetectShell(context.ShellName)
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
	// the test suite object for this file
	suite := junitxml.JUnitTestSuite{}
	suite.Name = inputfile
	suite.AddProperty("shelldoc-version", version.Version())
	// execute the interactions and verify the results:
	fmt.Printf("SHELLDOC: doc-testing \"%s\" ...\n", inputfile)
	results := resultStats{0, 0, 0}
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
		testcase := junitxml.JUnitTestCase{}
		fmt.Printf(opener, fmt.Sprintf("(%d)", index+1), interaction.Describe())
		testcase.Name = interaction.Cmd
		testcase.Classname = inputfile
		if context.Verbose {
			fmt.Printf(" --> %s\n", interaction.Cmd)
		}
		if err := interaction.Execute(&shell); err != nil {
			fmt.Printf(" --  ERROR: %v", err)
			context.RegisterReturnCode(returnError)
			results.errorCount++
		}
		fmt.Printf(closer, interaction.Result())
		if interaction.HasFailure() {
			context.RegisterReturnCode(returnFailure)
			suite.Failures++
			testcase.Failure = &junitxml.JUnitFailure{interaction.Result(), "failed", ""}
			results.failureCount++
		} else {
			results.successCount++
		}
		suite.RegisterTestCase(testcase)
		if interaction.HasFailure() && context.FailureStops {
			log.Printf("Stop requested after first failed test.")
			break
		}
	}
	fmt.Printf("%s: %d tests (%d successful, %d failures, %d execution errors)\n", result(context.ReturnCode()), suite.TestCount(), suite.SuccessCount(), results.failureCount, results.errorCount)
	context.Suites.Suites = append(context.Suites.Suites, suite)
	return results, nil
}

// ExecuteFiles runs each file through performInteractions and aggregates the results
func (context *Context) ExecuteFiles() int {
	context.RegisterReturnCode(returnSuccess)
	for _, file := range context.Files {
		_, err := context.performInteractions(file)
		if err != nil {
			fmt.Println(err) // log may be disabled (see "verbose")
			os.Exit(returnError)
		}
	}
	if err := context.WriteXML(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(returnError)
	}
	return context.ReturnCode()
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
