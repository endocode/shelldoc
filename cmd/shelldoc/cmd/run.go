// This file is part of shelldoc.
// © 2019, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: GPL-3.0

package cmd

import (
	"os"

	"github.com/endocode/shelldoc/pkg/run"
	"github.com/spf13/cobra"
)

var context run.Context

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a Markdown file as a documentation test",
	Long: `Run parses a Markdown input file, detects the code blocks in it,
executes them and compares their output with the content of the code block.`,
	Run: executeRun,
}

func init() {
	runCmd.Flags().StringVarP(&context.ShellName, "shell", "s", "", "The shell to invoke (default: $SHELL)")
	runCmd.Flags().BoolVarP(&context.FailureStops, "fail", "f", false, "Stop on the first failure")
	runCmd.Flags().StringVarP(&context.XMLOutputFile, "xml", "x", "", "Write results to the specified output file in JUnitXML format")
	runCmd.Flags().BoolVarP(&context.ReplaceDots, "replace-dots-in-xml-classname", "d", true, "When using filenames as classnames, replace dots with a unicode circle")
	rootCmd.AddCommand(runCmd)
}

func executeRun(cmd *cobra.Command, args []string) {
	context.Files = args
	os.Exit(context.ExecuteFiles())
}
