package junitxml

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func openTmpFile() (*os.File, error) {
	file, err := ioutil.TempFile("", "write_test-*.xml")
	if err != nil {
		return nil, fmt.Errorf("unable to open temporary output file: %v", err)
	}
	return file, nil
}

func removeTmpFile(filepath string) {
	const variable = "SHELLDOC_TEST_KEEP_TEMPORARY_FILES"
	if _, isSet := os.LookupEnv(variable); isSet == true {
		fmt.Printf("%s is set, not removing temporary file %s\n", variable, filepath)
		return
	}
	if err := os.Remove(filepath); err != nil {
		fmt.Fprintf(os.Stderr, "unable to remove temporary file at %s: %d", filepath, err)
	}
}

func validateXMLFile(filepath string) error {
	cmd := exec.Command("xmllint", "--noout", "--schema", "jenkins-junit.xsd", filepath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("XML validation finished with error: %v", err)
	}
	return nil
}

func TestMinimalDocument(t *testing.T) {
	// Write a minimal XML file with an empty testsuites section.
	testsuites := JUnitTestSuites{}

	file, err := openTmpFile()
	require.NoError(t, err, "Unable to open file for temporary XML document")
	defer removeTmpFile(file.Name())

	err = testsuites.Write(file)
	require.NoError(t, err, "Unable to write temporary XML document")
	// Verify it is schema compliant.
	require.NoError(t, validateXMLFile(file.Name()), "XML document fails to validate")
}

func TestOneTestSuite(t *testing.T) {
	// Write a minimal XML file with an empty testsuites section.
	testsuites := JUnitTestSuites{}
	ts := JUnitTestSuite{
		Tests:      1,
		Failures:   1,
		Time:       FormatTime(1234000000),
		Name:       "Test-TestSuite",
		Properties: []JUnitProperty{},
		TestCases:  []JUnitTestCase{},
	}
	ts.Properties = append(ts.Properties, JUnitProperty{"go.version", runtime.Version()})

	testCase := JUnitTestCase{
		Classname: "README.md",
		Name:      "ls -l",
		Time:      FormatTime(51345000),
		Failure: &JUnitFailure{
			Message:  "Failed",
			Type:     "mismatch",
			Contents: "(the test output)",
		},
	}
	ts.TestCases = append(ts.TestCases, testCase)
	testsuites.Suites = append(testsuites.Suites, ts)

	// The rest should be data/table driven...:
	file, err := openTmpFile()
	require.NoError(t, err, "Unable to open file for temporary XML document")
	defer removeTmpFile(file.Name())

	err = testsuites.Write(file)
	require.NoError(t, err, "Unable to write temporary XML document")
	// Verify it is schema compliant.
	require.NoError(t, validateXMLFile(file.Name()), "XML document fails to validate")
}
