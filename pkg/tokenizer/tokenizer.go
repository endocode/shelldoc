package tokenizer

// This file is part of shelldoc.
// © 2018, Mirko Boehm <mirko@endocode.com> and the shelldoc contributors
// SPDX-License-Identifier: LGPL-3.0

import (
	"log"
	"regexp"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

// Visitor contains the element handler functions
type Visitor struct {
	// CodeBlock should be assigned a function that will be called when a code block is encountered
	CodeBlock func(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus
	// FencedCodeBlock should be assigned a function to be called when a fenced code block is encountered
	FencedCodeBlock func(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus
	// After parsing, Interactions will hold the shell interactions found in the file
	Interactions []*Interaction
}

const cmdEx = "^[\\$>]\\s+(.+)$"

// handleCodeBlock parses the interactions in a code block and adds them to the Visitor
func handleCodeBlock(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus {
	cmdRx := regexp.MustCompile(cmdEx)

	lines := strings.Split(string(node.Literal), "\n")
	var current *Interaction
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		match := cmdRx.FindStringSubmatch(line)
		if len(match) > 1 {
			// begin a new command
			current = new(Interaction)
			visitor.Interactions = append(visitor.Interactions, current)
			cmd := match[1]
			current.Cmd = cmd
		} else {
			if current == nil {
				log.Printf("no trigger prefix ($ or >), skipping line: %s\n", line)
				continue
			}
			current.Response = append(current.Response, line)
		}
	}
	return blackfriday.GoToNext
}

// parseCodeBlockInfoString "best-faith" parses the info string and returns the language end the attributes
// if the info string is not written to the shelldoc specifications, both results are empty
func parseCodeBlockInfoString(infostring string) (string, map[string]string) {
	const infoStringHeaderEx = "^([.\\S]+)\\s+(.+)$"
	infoStringHeaderRx := regexp.MustCompile(infoStringHeaderEx)
	const attributesContentEx = "^.*\\{(.+)\\}.*$"
	attributesContentRx := regexp.MustCompile(attributesContentEx)
	const elementEx = "^([A-Za-z0-9]+)=(.+)$"
	elementRx := regexp.MustCompile(elementEx)

	var language string
	attributes := make(map[string]string)

	infostringmatch := infoStringHeaderRx.FindStringSubmatch(infostring)
	if infostringmatch != nil {
		language = infostringmatch[1]
		attributesString := infostringmatch[2]
		attributesContentMatch := attributesContentRx.FindStringSubmatch(attributesString)
		if attributesContentMatch != nil {
			attributesContent := attributesContentMatch[1]
			elements := strings.Split(attributesContent, " ")
			for _, element := range elements {
				if len(element) == 0 || !strings.HasPrefix(element, "shelldoc") {
					continue
				}
				elementmatch := elementRx.FindStringSubmatch(element)
				key := element
				value := ""
				if elementmatch != nil {
					key = elementmatch[1]
					value = elementmatch[2]
				}
				attributes[key] = value
			}
		} // else: ignore the rest of the infostring
	} // else: the info string is empty, treat this similar to a non-fenced code block

	return language, attributes
}

// handleFencedCodeBlock parses the interactions in a fenced code block and adds them to the Visitor
func handleFencedCodeBlock(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus {
	cmdRx := regexp.MustCompile(cmdEx)

	lines := strings.Split(string(node.Literal), "\n")
	if len(lines) < 2 {
		// technically, this should not happen, line 0 is the opening line of the code block (```),
		// the last line is the closer
		log.Printf("encountered a fenced code block with no info string, ignored")
		return blackfriday.GoToNext
	}
	infostring := lines[0]
	language, attributes := parseCodeBlockInfoString(infostring) // on error, language and attributes remain empty
	// closer := lines[len(lines)-1] // closer is not parsed any further
	lines = lines[1 : len(lines)-1]

	var current *Interaction
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		match := cmdRx.FindStringSubmatch(line)
		if len(match) > 1 {
			// begin a new command
			current = new(Interaction)
			current.Language = language
			current.Attributes = attributes
			visitor.Interactions = append(visitor.Interactions, current)
			cmd := match[1]
			current.Cmd = cmd
		} else {
			if current == nil {
				log.Printf("no trigger prefix ($ or >), skipping: %s\n", line)
				continue
			}
			current.Response = append(current.Response, line)
		}
	}
	return blackfriday.GoToNext
}

// NewInteractionVisitor creates a visitor configured with the default ineraction parser
func NewInteractionVisitor() *Visitor {
	visitor := new(Visitor)
	visitor.CodeBlock = handleCodeBlock
	visitor.FencedCodeBlock = handleFencedCodeBlock
	return visitor
}

// visit is called on every Markdown element encountered
// It checks for code blocks and calls the respective handlers.
func (visitor *Visitor) visit(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	// log.Printf("%v: %s", node.Type, node.Literal)
	if node.Type == blackfriday.CodeBlock && entering == true {
		return visitor.CodeBlock(visitor, node)
	} else if node.Type == blackfriday.Code && entering == true {
		return visitor.FencedCodeBlock(visitor, node)
	}
	return blackfriday.GoToNext
}

// Tokenize parses the data and calls the event handlers on visitor
func Tokenize(data []byte, visitor *Visitor) error {
	md := blackfriday.New()
	om := md.Parse(data)
	om.Walk(visitor.visit)
	return nil
}
