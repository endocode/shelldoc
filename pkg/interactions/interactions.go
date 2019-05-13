package interactions

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/endocode/shelldoc/pkg/junitxml"
	"github.com/endocode/shelldoc/pkg/shell"
	"github.com/endocode/shelldoc/pkg/tokenizer"
)

func max(a, b int) int { // really, golang?
	if a > b {
		return a
	}
	return b
}

type resultStats struct {
	returncode, testCount, successCount, failureCount, errorCount int
}

const (
	returnSuccess = iota
	returnFailure
	returnError
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

func performInteractions(inputfile string, shellname string, verbose bool, failureStops bool, suites *junitxml.JUnitTestSuites) (resultStats, error) {
	// detect shell
	shellpath, err := shell.DetectShell(shellname)
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
	results := resultStats{returnSuccess, 0, 0, 0, 0}
	// construct the opener and closer format strings, since they depend on verbose mode
	magnitude := int(math.Log10(float64(len(visitor.Interactions)))) + 1
	openerLineEnding := "  : "
	resultString := " "
	if verbose {
		openerLineEnding = "\n"
		resultString = " <-- "
	}
	counterFormat := fmt.Sprintf("%%%ds", magnitude+2)
	opener := fmt.Sprintf(" CMD %s: %%s%s", counterFormat, openerLineEnding)
	closer := fmt.Sprintf("%s%%s\n", resultString)

	for index, interaction := range visitor.Interactions {
		testcase := junitxml.JUnitTestCase{}
		suite.Tests++
		results.testCount++
		fmt.Printf(opener, fmt.Sprintf("(%d)", index+1), interaction.Describe())
		testcase.Name = interaction.Cmd
		testcase.Classname = inputfile
		if verbose {
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
			suite.Failures++
			testcase.Failure = &junitxml.JUnitFailure{interaction.Result(), "failed", ""}
			results.failureCount++
		} else {
			results.successCount++
		}
		suite.TestCases = append(suite.TestCases, testcase)
		if interaction.HasFailure() && failureStops {
			log.Printf("Stop requested after first failed test.")
			break
		}
	}
	fmt.Printf("%s: %d tests (%d successful, %d failures, %d execution errors)\n", result(results.returncode), results.testCount, results.successCount, results.failureCount, results.errorCount)
	suites.Suites = append(suites.Suites, suite)
	return results, nil
}

// ExecuteFiles runs each file through performInteractions and aggregates the results
func ExecuteFiles(files []string, shellname string, verbose bool, failureStops bool) int {
	returnCode := returnSuccess
	suites := junitxml.JUnitTestSuites{}
	for _, file := range files {
		results, err := performInteractions(file, shellname, verbose, failureStops, &suites)
		if err != nil {
			fmt.Println(err) // log may be disabled (see "verbose")
			os.Exit(returnError)
		}
		returnCode = max(results.returncode, returnCode)
	}
	// write the result to the specified XML output file:
	writeXML := true
	xmlFilename := "testresults.xml"
	if writeXML {
		file, err := os.OpenFile(xmlFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			fmt.Printf("Unable to open XML output file for writing: %v\n", err)
			os.Exit(returnError)
		}
		if err := suites.Write(file); err != nil {
			fmt.Printf("Error writing XML output file: %v\n", err)
		}
	}
	return returnCode
}
