package tokenizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Endocode/shelldoc/pkg/interaction"
	"gopkg.in/russross/blackfriday.v2"
)

// Visitor contains the element handler functions
type Visitor struct {
	// Assign a function that will be called when a code block is encountered
	CodeBlock func(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus
	// After parsing, Interactions will hold the shell interactions found in the file
	Interactions []*interaction.Interaction
}

// ParseInteractions parses the interactions in a code block and adds them to the Visitor
func ParseInteractions(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus {
	cmdEx := "^[\\$>]\\s+(.+)$"
	cmdRx := regexp.MustCompile(cmdEx)

	lines := strings.Split(string(node.Literal), "\n")
	var current *interaction.Interaction
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		//fmt.Printf("%4d: %s\n", counter, line)
		match := cmdRx.FindStringSubmatch(line)
		if len(match) > 1 {
			// begin a new command
			current = new(interaction.Interaction)
			visitor.Interactions = append(visitor.Interactions, current)
			cmd := match[1]
			current.Cmd = cmd
		} else {
			if current == nil {
				fmt.Printf("Skipping line since there was no command: %s", line)
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
func Tokenize(data []byte, visitor *Visitor) error {
	md := blackfriday.New()
	om := md.Parse(data)
	om.Walk(visitor.visit)
	return nil
}
