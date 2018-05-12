package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ReadInput reads either the files specified on the command line or stdin and returns the bytes.
// Markdown.Parse expects bytes, not a stream.
func ReadInput(args []string) ([]byte, error) {
	if len(args) > 0 {
		var result []byte
		for _, filename := range args {
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, fmt.Errorf("unable to read file %s", filename)
			}
			result = append(result, content[:]...)
		}
		return result, nil
	} else {
		result, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			if err != nil {
				return nil, fmt.Errorf("unable to read from stdin: %v", err)
			}
		}
		return result, nil
	}
}
