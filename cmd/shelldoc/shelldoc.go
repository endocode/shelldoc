package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Endocode/shelldoc/pkg/shell"
	"github.com/Endocode/shelldoc/pkg/tokenizer"
)

const (
	returnSuccess = iota
	returnFailure
	returnError
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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

func main() {
	data, err := ReadInput(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
	// run the input through the tokenizer
	visitor := tokenizer.NewInteractionVisitor()
	tokenizer.Tokenize(data, visitor)
	// start a background shell, it will run until the program ends
	shell, err := shell.StartShell()
	if err != nil {
		log.Fatalln(err)
	}
	defer shell.Exit()
	returncode := returnSuccess
	testCount := 0
	successCount := 0
	failureCount := 0
	errorCount := 0
	// execute the interactions and verify the results:
	for index, interaction := range visitor.Interactions {
		testCount++
		fmt.Printf("COMMAND (%d): %s\n", index+1, interaction.Describe())
		fmt.Printf("--> %s\n", interaction.Cmd)
		if err := interaction.Execute(&shell); err != nil {
			fmt.Printf("--  ERROR: %v\n", err)
			returncode = max(returncode, returnError)
			errorCount++
		}
		fmt.Printf("<-- %s\n", interaction.Result())
		if interaction.HasFailure() {
			returncode = max(returncode, returnFailure)
			failureCount++
		} else {
			successCount++
		}
	}
	fmt.Printf("%s: %d tests (%d successful, %d failures, %d execution errors)\n", result(returncode), testCount, successCount, failureCount, errorCount)
	os.Exit(returncode)
}
