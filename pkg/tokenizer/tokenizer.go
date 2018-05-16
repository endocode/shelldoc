package tokenizer

import (
	"fmt"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

// Visitor contains the element handler functions
type Visitor struct {
	// Assign a function that will be called when a code block is encountered
	CodeBlock func(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus
	// After parsing, Interactions will hold the shell interactions found in the file
	Interactions []Interaction
}

// Interaction represents one interaction with the shell
type Interaction struct {
	Cmd              string
	Response         string
	AlternativeRegEx string
}

// ParseInteractions parses the interactions in a code block and adds them to the Visitor
func ParseInteractions(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus {
	lines := strings.Split(string(node.Literal), "\n")
	fmt.Printf("node: %s - %s.\n", node.Type, lines)
	return blackfriday.GoToNext
}

// NewInteractionVisitor creates a visitor configured with the default ineraction parser
func NewInteractionVisitor() *Visitor {
	visitor := new(Visitor)
	visitor.CodeBlock = ParseInteractions
	return visitor
}

// visit is called on every Markdown element encountered
// It checks for code blocks and calls the respective handlers.
func (visitor *Visitor) visit(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	//fmt.Printf("node: %s - %s - entering: %v.\n", node.Type, node.Literal, entering)
	if node.Type == blackfriday.CodeBlock && entering == true {
		return visitor.CodeBlock(visitor, node)
	}
	return blackfriday.GoToNext
}

// Tokenize parses the data and calls the event handlers on visitor
func Tokenize(data []byte, visitor Visitor) error {
	md := blackfriday.New()
	om := md.Parse(data)
	om.Walk(visitor.visit)
	return nil
}
