// This file is part of shelldoc.
// © 2019, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: GPL-3.0

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shelldoc",
	Short: "shelldoc tests Unix shell commands in Markdown documentation",
	Long: `Markdown is widely used for documentation and README.md files that explain how
to use or build some software. Such documentation often contains shell commands that
explain how to build a software or how to run it. To make sure the documentation is a
ccurate and up-to-date, it should be automatically tested. shelldoc tests Unix shell
commands in Markdown files and reports the results.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogging)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable diagnostics")
}

func initLogging() {
	// verbose essentially enables or disables log output:
	if verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(0)
	log.SetPrefix("Note: ")
}
