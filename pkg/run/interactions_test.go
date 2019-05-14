package run

// This file is part of shelldoc.
// © 2018, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: Apache-2.0

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	verbose      bool
	shellname    string
	failureStops bool
)

func TestMain(m *testing.M) {
	verbose = true
	failureStops = false
	os.Exit(m.Run())
}
func TestHelloWorld(t *testing.T) {
	context := Context{
		// ShellName:    "",
		Verbose: true,
		// FailureStops: false,
		// XMLOutputFile: "",
		// Files:         []string{},
		// Suites: junitxml.JUnitTestSuites{},
	}
	results, err := context.performInteractions("../../pkg/tokenizer/samples/helloworld.md")
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnSuccess, context.ReturnCode(), "The expected return code is returnSuccess.")
	require.Equal(t, 4, results.successCount, "There are three successful tests in the sample.")
}

func TestHFailNoMatch(t *testing.T) {
	context := Context{}
	results, err := context.performInteractions("../../pkg/tokenizer/samples/failnomatch.md")
	require.NoError(t, err, "The failnomatch example should fail with a mismatch.")
	require.Equal(t, returnFailure, context.ReturnCode(), "The expected return code is returnFailure.")
	require.Equal(t, 1, results.failureCount, "There is one failing test in the sample.")
}

func TestExitCodesOptions(t *testing.T) {
	context := Context{}
	_, err := context.performInteractions("../../pkg/tokenizer/samples/options.md")
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnSuccess, context.ReturnCode(), "The expected return code is returnSuccess.")
}
