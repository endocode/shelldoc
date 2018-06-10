package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	options.verbose = true
	initializeLogging()
	os.Exit(m.Run())
}
func TestHelloWorld(t *testing.T) {
	results, err := performInteractions("../../pkg/tokenizer/samples/helloworld.md")
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnSuccess, results.returncode, "The expected return code is returnSuccess.")
	require.Equal(t, 3, results.successCount, "There are three successful tests in the sample.")
}

func TestHFailNoMatch(t *testing.T) {
	results, err := performInteractions("../../pkg/tokenizer/samples/failnomatch.md")
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnFailure, results.returncode, "The expected return code is returnFailure.")
	require.Equal(t, 1, results.failureCount, "There is one failing test in the sample.")
}
