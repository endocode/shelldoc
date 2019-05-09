// This file is part of shelldoc.
// © 2019, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: GPL-3.0

package cmd

import (
	"os"

	"github.com/endocode/shelldoc/pkg/interactions"
	"github.com/spf13/cobra"
)

var shellname string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a Markdown file as a documentation test",
	Long: `Run parses a Markdown input file, detects the code blocks in it,
executes them and compares their output with the content of the code block.`,
	Run: executeRun,
}

func init() {
	runCmd.Flags().StringVarP(&shellname, "shell", "s", "", "Keep the created directories instead of cleaning up.")
	rootCmd.AddCommand(runCmd)
}

func executeRun(cmd *cobra.Command, args []string) {
	os.Exit(interactions.ExecuteFiles(args, shellname, verbose))
}
