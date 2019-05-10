package junitxml

import (
	"encoding/xml"
	"fmt"
	"io"
)

func (testsuites JUnitTestSuites) Write(w io.Writer) error {
	io.WriteString(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	enc := xml.NewEncoder(w)
	enc.Indent("", "\t")
	if err := enc.Encode(testsuites); err != nil {
		fmt.Printf("unable to write XML document: %v", err)
	}
	io.WriteString(w, "\n")
	return nil
}
