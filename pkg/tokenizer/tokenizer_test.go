package tokenizer

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var echoTrueCodeBlockCount int

func codeBlockHandler(node *blackfriday.Node) blackfriday.WalkStatus {
	// fmt.Printf("node: %s - %s.\n", node.Type, node.Literal)
	echoTrueCodeBlockCount++
	return blackfriday.GoToNext
}
func TestEchoTrue(t *testing.T) {
	data, err := ioutil.ReadFile("samples/echotrue.md")
	require.NoError(t, err, "Unable to read sample data file")
	visitor := Visitor{codeBlockHandler}
	require.Zero(t, echoTrueCodeBlockCount, "Starting the counter")
	Tokenize(data, visitor)
	require.Equal(t, echoTrueCodeBlockCount, 1, "There is one code block element in the sample file")
}
