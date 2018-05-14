package tokenizer

import (
	"gopkg.in/russross/blackfriday.v2"
)

// Visitor contains the element handler functions
type Visitor struct {
	// Assign a function that will be called when a code block is encountered
	CodeBlock func(node *blackfriday.Node) blackfriday.WalkStatus
}

func (visitor *Visitor) visit(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	//fmt.Printf("node: %s - %s - entering: %v.\n", node.Type, node.Literal, entering)
	if node.Type == blackfriday.CodeBlock {
		return visitor.CodeBlock(node)
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
