package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Endocode/shelldoc/pkg/shell"
)

func TestHelloWorld(t *testing.T) {
	options.verbose = true
	initializeLogging()
	shell, err := shell.StartShell()
	if err != nil {
		Log.Fatalln(err)
	}
	defer shell.Exit()

	results, err := performInteractions([]string{"../../pkg/tokenizer/samples/helloworld.md"}, &shell)
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnSuccess, results.returncode, "The expected return code is returnSuccess.")
	require.Equal(t, 3, results.successCount, "There are three successful tests in the sample.")
}

func TestHFailNoMatch(t *testing.T) {
	options.verbose = true
	initializeLogging()
	shell, err := shell.StartShell()
	if err != nil {
		Log.Fatalln(err)
	}
	defer shell.Exit()

	results, err := performInteractions([]string{"../../pkg/tokenizer/samples/failnomatch.md"}, &shell)
	require.NoError(t, err, "The HelloWorld example should execute without errors.")
	require.Equal(t, returnFailure, results.returncode, "The expected return code is returnFailure.")
	require.Equal(t, 1, results.failureCount, "There is one failing test in the sample.")
}
