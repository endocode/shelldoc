package interactions

// This file is part of shelldoc.
// © 2018, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: Apache-2.0

import (
	"os"
	"testing"

	"github.com/endocode/shelldoc/pkg/junitxml"
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
	results, err := performInteractions("../../pkg/tokenizer/samples/helloworld.md", "", verbose, failureStops, &junitxml.JUnitTestSuites{})
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnSuccess, results.returncode, "The expected return code is returnSuccess.")
	require.Equal(t, 4, results.successCount, "There are three successful tests in the sample.")
}

func TestHFailNoMatch(t *testing.T) {
	results, err := performInteractions("../../pkg/tokenizer/samples/failnomatch.md", "", verbose, failureStops, &junitxml.JUnitTestSuites{})
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnFailure, results.returncode, "The expected return code is returnFailure.")
	require.Equal(t, 1, results.failureCount, "There is one failing test in the sample.")
}

func TestExitCodesOptions(t *testing.T) {
	results, err := performInteractions("../../pkg/tokenizer/samples/options.md", "", verbose, failureStops, &junitxml.JUnitTestSuites{})
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnSuccess, results.returncode, "The expected return code is returnFailure.")
}
