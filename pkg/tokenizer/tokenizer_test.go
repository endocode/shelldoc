package tokenizer

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var echoTrueCodeBlockCount int

func codeBlockHandler(visitor *Visitor, node *blackfriday.Node) blackfriday.WalkStatus {
	echoTrueCodeBlockCount++
	return blackfriday.GoToNext
}
func TestEchoTrue(t *testing.T) {
	data, err := ioutil.ReadFile("samples/echotrue.md")
	require.NoError(t, err, "Unable to read sample data file")
	visitor := Visitor{codeBlockHandler, nil}
	require.Zero(t, echoTrueCodeBlockCount, "Starting the counter")
	Tokenize(data, &visitor)
	require.Equal(t, echoTrueCodeBlockCount, 1, "There is one code block element in the sample file")
}

func TestTokenizeEchoTrue(t *testing.T) {
	data, err := ioutil.ReadFile("samples/echotrue.md")
	require.NoError(t, err, "Unable to read sample data file")
	visitor := NewInteractionVisitor()
	Tokenize(data, visitor)
	require.Equal(t, len(visitor.Interactions), 1, "There is one code block element in the sample file")
}
