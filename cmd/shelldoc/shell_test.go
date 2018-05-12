package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellLifeCycle(t *testing.T) {
	// The most basic test, start a shell and exit it again
	shell, err := startShell()
	require.NoError(t, err, "Starting a shell should work")
	require.NoError(t, shell.exit(), "Exiting ad running shell should work")
}

func TestShellLifeCycleRepeated(t *testing.T) {
	// Can the program start and stop a shell repeatedly?
	for counter := 0; counter < 16; counter++ {
		shell, err := startShell()
		require.NoError(t, err, "Starting a shell should work")
		require.NoError(t, shell.exit(), "Exiting ad running shell should work")
	}
}

func TestReturnCodes(t *testing.T) {
	// Does the shell report return codes corrrectly?
	shell, err := startShell()
	require.NoError(t, err, "Starting a shell should work")
	defer shell.exit()
	{
		output, rc, err := shell.executeCommand("true")
		require.NoError(t, err, "The true command is a builtin and should always work")
		require.Equal(t, 0, rc, "The exit code of true should always be zero")
		require.Empty(t, output, "true does not say a word")
	}
	{
		output, rc, err := shell.executeCommand("false")
		require.NoError(t, err, "The false command is a builtin and should always work")
		require.NotEqual(t, 0, rc, "The exit code of false should never be zero")
		require.Empty(t, output, "false does not say a word")
	}
}

func TestCaptureOutput(t *testing.T) {
	// Does the shell capture and return the lines printed by the command correctly?
	shell, err := startShell()
	require.NoError(t, err, "Starting a shell should work")
	defer shell.exit()
	{
		const (
			hello = "Hello"
			world = "World"
		)
		output, rc, err := shell.executeCommand(fmt.Sprintf("echo %s && echo %s", hello, world))
		require.NoError(t, err, "The echo command is a builtin and should always work")
		require.Equal(t, 0, rc, "The exit code of echo should be zero")
		require.Len(t, output, 2, "echo was called twice")
		require.Equal(t, output[0], hello, "you had one job, echo")
		require.Equal(t, output[1], world, "actually, two")
	}
}
